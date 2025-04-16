# SO_examen
Instrucciones:  
Primera Parte:  
Soy un ingeniero de Devops que está encargado de integrar un
pipeline de CI/CD a un proyecto en Godot, tengo armado un
pequeño esqueleto de un contenedor base para compilar los
proyectos. Tu deber es completar el contenedor para ejecutar
Godot en modo sin interfaz.
Usen el documento adjunto al examen.  
Ejemplo de output al ejecutar el contenedor:  
```
Godot Engine v4.4.stable.official.4c311cbee - https://godotengine.org
```
Segunda Parte:
Integren al API con Go una vista creada con React, la vista
dependerá de qué tipo de API realizaron para practicar.
El API de Go puede estar conectado a la base de datos de su
elección.
Tienen que orquestar la ejecución de los servicios por medio
de contenedores con docker compose, teniendo el servicio de
API y de Frontend con react una configuración personalizada en
sus respectivos Dockerfiles.  
EXTRA:  
Con su contenedor de godot compilen de manera exitosa el
proyecto de https://godotengine.org/asset-library/asset/540

## Tabla de contenido
1. [Introduccion](#introduccion)  
2. [Parte 1](#Parte-1)  
3. [Parte 2](#Parte-2)  
    a. [MySQL](#Contenedor-MySQL)  
    b. [Go](#Contenedor-Go)  
    c. [Node](#Contenedor-Node)  

## Introduccion
El objetivo principal es ejecutar y poner en practica lo aprendido sobre Docker.  
Para la primera parte se trata de completar el Dockerfile ejecutando Godot de manera correcta.  
La segunda parte consta de desarrollar 3 contendores los cuales se conecten en cadena. El primero contedor por eleccion personal usa MySQL como motor de base de datos, el segundo contenedor contiene Golang en el cual se desarrollo una API conectada a la base de datos del primer contenedor. El tercer contenedor contiene Node, el cual ejecuta una aplicacion de React el cual obtiene los datos de la API y los muestra en pantalla.

## Parte 1
El contenedor con el cual empezamos es el siguiente:  
```
FROM ubuntu:24.04

RUN wget https://github.com/godotengine/godot/releases/download/4.4-stable/Godot_v4.4-stable_linux.x86_64.zip -O /tmp/godot.zip \
&& rm /tmp/godot.zip \
&& mv /usr/local/bin/Godot_v4.4-stable_linux.x86_64 /usr/local/bin/godot+

CMD ["godot", "--headless"]
```
Si tratamos de ejecutar lo siguiente genera un error.

## Parte 2
El desarrollo de la parte 2 tomo mucho tiempo, pues se soluciona un error pero despues se genera otro. Es importante resaltar que implico mucha investigacion para lograr el siguiente proyecto ya que nunca habia trabajado con React a excepcion de ejecutar un proyecto y desarrollar una API en Go es mucho mas avanzado que simplemente programar algo basico en Go.  
Se ha seleccionado MySQL como motor para la base de datos ya que es con el que mas experiencia tengo ya que las materias de bases de dato se han desarrollado con este motor.  

### Contenedor MySQL
Para el contenedor MySQL, se usa la version 8.0 de MySQL, porque es la version que uso con el sistema de la computadora.  
La base de datos es la base de datos 'classicmodels' la cual es la que hemos usado en la clase de bases de datos.

El Dockerfile para el contenedor es el siguiente: 
```
FROM mysql:8.0
COPY ./mysqlsampledatabase.sql /docker-entrypoint-initdb.d/
ENV MYSQL_ROOT_PASSWORD=root
```
Lo realmente importante es que se copia el documento con el dump de la base de datos en el directorio '/docker-entrypoint-initdb.d/' el cual se ejecuta siempre que se ejecuta el contenedor por lo que una vez ejecutado el contenedor ya tiene la base de datos en mysql.  


### Contenedor Go

Para Go tenemos 4 documentos:  
    Dockerfile  
    go.mod  
    go.sum  
    main.go  

go.mod: es el documento de modulo de Go. Aqui es donde se anotan cuales son las dependencias del modulo en cuestion. Se usa para poder usar las dependencias necesarias externas en el programa o paquetes que se usen en un programa. Este documento se modifica solo usando el comando para descargar dependencias externas.  

go.sum: se modifica al igual que con el go.mod y este tiene la lista de las dependencias almacenas y es necesario que haya concordancia entre unos y otros.

#### main.go

Dependencias:  
"github.com/go-sql-driver/mysql"  
Es el driver que permite la conexion entre go y mysql. Necesario para traer los datos de la base de datos a la API  

"github.com/gin-gonic/gin"  
Es un Framework Web que sirve para desarrollar la API en Go. Con este se pueden definir los endpoint, asi como definir las funciones por las cuales se transmitiran los datos.  

"github.com/gin-contrib/cors"  
Este modulo sirve para el problema con CORS en la API, en caso de no generar ni usar modificacion alguna no permite acceso de otras aplicaciones. Por lo tanto se usa para definir los permisos de acceso.


### Contenedor Node
