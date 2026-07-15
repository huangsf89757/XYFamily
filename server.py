#!/usr/bin/env python3
"""
XYFamily Wiki 开发服务器
提供静态文件服务，Markdown 文件由前端 wiki-browser.html 内置渲染器处理。
"""

import http.server
import socketserver
import os
from urllib.parse import unquote, quote

# 浏览页目录（相对于项目根）
BROWSE_DIR = os.path.join('wiki', '浏览页')


class WikiRequestHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        path = unquote(self.path)

        # 根路径 → 302 重定向到浏览页目录（保持相对链接正确）
        if path in ('/', '/index.html', '/wiki', '/wiki/'):
            self.send_response(302)
            self.send_header('Location', quote('/wiki/浏览页/index.html'))
            self.end_headers()
            return

        # .md 文件 → 返回原始 Markdown 文本（由前端渲染器处理）
        if path.endswith('.md'):
            filepath = '.' + path
            if os.path.exists(filepath):
                self.send_response(200)
                self.send_header('Content-type', 'text/markdown; charset=utf-8')
                self.send_header('Cache-Control', 'no-cache')
                self.end_headers()
                with open(filepath, 'rb') as f:
                    self.wfile.write(f.read())
                return
            else:
                self.send_error(404, f'文件未找到: {path}')
                return

        # 其他静态文件交给父类处理
        return super().do_GET()

    def end_headers(self):
        # 添加 CORS 头，允许前端 fetch 跨域
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET')
        super().end_headers()

    def guess_type(self, path):
        """确保 .md 文件返回 text/markdown MIME 类型"""
        if path.endswith('.md'):
            return 'text/markdown; charset=utf-8'
        return super().guess_type(path)


def main():
    PORT = 8000

    for port in [8000, 8001, 8002, 8003, 8004]:
        try:
            with socketserver.TCPServer(("", port), WikiRequestHandler) as httpd:
                print(f"🎉 XYFamily Wiki 开发服务器已启动!")
                print(f"🔗 预览地址: http://localhost:{port}")
                print(f"📁 文档根目录: {os.getcwd()}")
                print(f"📄 首页: http://localhost:{port}/")
                print(f"🚀 按 Ctrl+C 停止服务器")
                print("=" * 50)
                print("✨ 特性:")
                print("  • 浏览页前端内置 Markdown 渲染器")
                print("  • .md 文件返回原始文本，由前端渲染")
                print("  • 静态文件服务 + CORS 支持")
                print("=" * 50)

                httpd.serve_forever()
                break
        except OSError as e:
            if port == 8004:
                print(f"❌ 无法启动服务器，端口8000-8004都被占用了")
                print(f"   错误信息: {e}")
                exit(1)
            else:
                print(f"⚠️  端口 {port} 被占用，尝试端口 {port+1}...")


if __name__ == "__main__":
    main()
