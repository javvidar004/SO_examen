FROM ubuntu:24.04

# Instala unzip y wget
RUN apt-get update && apt-get install -y wget unzip

# Descarga el archivo ZIP de Godot
RUN wget https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip -O /tmp/godot.zip

# Crea un directorio para Godot y extrae el ejecutable
RUN mkdir -p /usr/local/bin/godot && \
    unzip /tmp/godot.zip -d /usr/local/bin/godot && \
    rm /tmp/godot.zip

# Renombra el ejecutable y le da permisos de ejecución
RUN mv /usr/local/bin/godot/Godot_v4.4-stable_linux.x86_64 /usr/local/bin/godot/godot && \
    chmod +x /usr/local/bin/godot/godot

# Establece el ejecutable de Godot como el comando principal
CMD ["/usr/local/bin/godot/godot", "--headless"]
