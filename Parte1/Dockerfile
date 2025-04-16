FROM ubuntu:24.04

RUN apt-get update && apt-get install -y wget unzip fontconfig

RUN wget https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip -O /tmp/godot.zip \
&& unzip /tmp/godot.zip \
&& rm /tmp/godot.zip \
&& mv Godot_v4.4-stable_linux.x86_64 /usr/local/bin/godot

WORKDIR /usr/local/bin

CMD ["godot", "--headless"]
