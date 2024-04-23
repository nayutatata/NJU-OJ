# 使用Ubuntu作为基础镜像  
FROM real_2

# 设置工作目录  
WORKDIR /app  


# 设置环境变量  
ENV PATH="/home/dockeruser/.local/bin:${PATH}"  

# 复制你的C++源代码和Python判题脚本到镜像中  
COPY . /app  

USER root

# CMD [ "/bin/bash" ]
# 编译你的C++代码（如果需要）  
# RUN g++ -o my_program my_program.cpp  

# 暴露端口（如果需要运行网络服务）  
EXPOSE 9001

# 设置默认命令（根据你的应用设置）  
CMD ["/usr/bin/tail", "-f", "/dev/null"]  

# 如果你的应用是一个服务，你可能需要运行它  
# CMD ["your-service-command"]