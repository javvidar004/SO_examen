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

Despues se definen las constantes las cuales son los datos de acceso para la base de datos en cuestion. Para el caso de dbHost no se define la IP del contenedor pues este tiende a cambiar y por la red de Docker puede cambiar. Al docker tener integrado DNS se puede usar con el nombre que se define del contenedor en el docker-compose.  

Se define la estructura sobre la cual se recibiran los datos de la base de datos y despues se enviaran en formato JSON al frontend.  

Despues se encuentra la funcion getDBConfig() esta funcion usa los datos de conexion para crear la DSN (Data Source Name) el cual es la linea de texto, que describe la conexion con la base de datos. Ademas al darle formato DSN no es necesario estructurar la linea de texto directamente y es mas claro el acceso y la declaracion de los datos.  

```
getDBConfig()
```

Seguido esta la funcion getItems() esta es la funcion que genera el request a la base de datos y tambien mas adelante se declara como la funcion a ejecutar en un determinado endpoint de la API. Por ser usada como funcion que manda datos en formato JSON, pasa por parametro el contexto de Gin, que es lo que enviara al frontend.  
Dentro de la misma funcion se genera la conexion a la base de datos, se ejecuta el query definido, se almacenan los datos en un slice para organizarlos y darles formato.  
Se manejan las diferentes excepciones ya sea por si el query fallo, el escaneo de datos, o la iteracion al manejar los datos.
Finalmente si todo esta bien, regresa los datos en formato JSON y con codigo 200.  
```
getItems()
```
La siguiente es importante para saber en caso de errores donde se ha ocasionado el error dentro de la API.  
```
errorHandler()
```

Finalmente en el main se define el funcionamiento de Gin.  
Despues la configuracion de CORS para permiter el acceso y las solicitudes desde el frontend.  
Se introduce el manejador de errores a Gin.  
Se define el endpoint de la API.  
Y finalmente se ejecuta el servidor en el puerto 8080.  

Por ultimo, el Dockerfile copia los archivos de las dependencias para despues ejecutar la instalacion de los mismos, despues copia el programa de la API. Para finalmente declarar el comando para que la API este en ejecucion.


### Contenedor Node

Para este contenedor primero se creo la aplicacion en React. El unico documento modificado fue App.js y se agrego el documento TablaDatos.js. En App.js solo se agrega el componente que se declara en TablaDatos.js.  

En tabla datos primero se declara la URL de la api. Para declarar la URL primero tenia un error el cual era utilizar la URL como si todo se estuviera ejecutando desde el lado del servidor porque ponia el URL 'http://api:8080/items' y generaba error haciendome pensar que el DNS de Docker estaba fallando pero resulta que al estarse ejecutando del lado del cliente habia que hacer la llamada desde el localhost, por lo que la URL definitiva fue 'http://localhost:8080/items'.  

Dentro solo se define una funcion que es con la cual se reciben los datos y se muestran en el HTML. Se manejan excepciones en la obtencion de los datos, pero en caso de obtenerlos se pasa a la parte inferior en la cual se define la tabla con en la cual se mostraran los datos obtenidos de la API.  

Los documentos package.json y package-lock.json son los documentos de las dependencias de para ejecutar la aplicacion. Los cuales son esenciales para la ejecucion de la aplicacion en el contenedor.  

El Dockerfile copia los datos que estan en la carpeta de la aplicacion y despues ejecuta
```
npm install
```
El cual instala las dependencias para despues definir, la ejecucion del programa y que siempre que se ejecute el contenedor se ejecute la aplicacion.  

