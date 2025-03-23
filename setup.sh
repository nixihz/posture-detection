#!/bin/bash

# OpenCV 安装路径
OPENCV_PATH="/opt/homebrew/Cellar/opencv/4.11.0_1"

# 设置 OpenCV 的 pkg-config 路径
export PKG_CONFIG_PATH="$OPENCV_PATH/lib/pkgconfig:$PKG_CONFIG_PATH"

# 设置 OpenCV 的库路径
export LIBRARY_PATH="$OPENCV_PATH/lib:$LIBRARY_PATH"

# 设置 OpenCV 的包含路径
export CGO_CPPFLAGS="-I$OPENCV_PATH/include"
export CGO_LDFLAGS="-L$OPENCV_PATH/lib -lopencv_core -lopencv_imgproc -lopencv_imgcodecs -lopencv_videoio -lopencv_objdetect"

# 运行程序
go run cmd/main.go