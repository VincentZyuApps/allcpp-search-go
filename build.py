#!/usr/bin/env python3
"""
CPP Search Go 编译脚本
用法: python build.py --version v1.0.0
"""

import argparse
import subprocess
import os
import sys
from datetime import datetime

# 编译目标平台
TARGETS = [
    {"goos": "windows", "goarch": "amd64", "ext": ".exe"},
    {"goos": "linux", "goarch": "amd64", "ext": ""},
]

# 项目名称
PROJECT_NAME = "cpp_search"

def run_command(cmd: list, env: dict = None) -> bool:
    """运行命令"""
    full_env = os.environ.copy()
    if env:
        full_env.update(env)
    
    print(f"  → {' '.join(cmd)}")
    result = subprocess.run(cmd, env=full_env, capture_output=True, text=True)
    
    if result.returncode != 0:
        print(f"  ✗ 错误: {result.stderr}")
        return False
    return True

def build(version: str, output_dir: str = "bin"):
    """编译所有目标平台"""
    print(f"🚀 CPP Search Go 编译脚本")
    print(f"📌 版本: {version}")
    print(f"📅 时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("-" * 50)
    
    # 创建输出目录
    os.makedirs(output_dir, exist_ok=True)
    
    success_count = 0
    failed_count = 0
    
    for target in TARGETS:
        goos = target["goos"]
        goarch = target["goarch"]
        ext = target["ext"]
        
        output_name = f"{PROJECT_NAME}_{goos}_{goarch}_{version}{ext}"
        output_path = os.path.join(output_dir, output_name)
        
        print(f"\n🔨 编译 {goos}/{goarch}...")
        
        env = {
            "GOOS": goos,
            "GOARCH": goarch,
            "CGO_ENABLED": "0",  # 禁用 CGO，生成静态链接二进制
        }
        
        # 使用 ldflags 注入版本信息
        ldflags = f"-s -w -X main.Version={version}"
        cmd = ["go", "build", "-ldflags", ldflags, "-o", output_path, "."]
        
        if run_command(cmd, env):
            # 获取文件大小
            size = os.path.getsize(output_path)
            size_mb = size / (1024 * 1024)
            print(f"  ✓ 成功: {output_name} ({size_mb:.2f} MB)")
            success_count += 1
        else:
            print(f"  ✗ 失败: {output_name}")
            failed_count += 1
    
    print("\n" + "=" * 50)
    print(f"📊 编译完成: {success_count} 成功, {failed_count} 失败")
    print(f"📁 输出目录: {os.path.abspath(output_dir)}")
    
    if failed_count > 0:
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="CPP Search Go 编译脚本")
    parser.add_argument(
        "--version", "-v",
        required=True,
        help="版本号，例如: v1.0.0"
    )
    parser.add_argument(
        "--output", "-o",
        default="bin",
        help="输出目录 (默认: bin)"
    )
    
    args = parser.parse_args()
    build(args.version, args.output)

if __name__ == "__main__":
    main()
