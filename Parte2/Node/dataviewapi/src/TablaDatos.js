import React, { useState, useEffect } from 'react';

// URL de tu API (asegúrate de que sea accesible)
// Puede que necesites http:// o https://
const apiUrl = 'http://localhost:8080/items'; // <-- AJUSTA ESTO SI ES NECESARIO

function TablaDatos() {
    // Estado para almacenar los datos de la API
    const [data, setData] = useState([]); // Inicializa como array vacío
    // Estado para manejar el estado de carga
    const [loading, setLoading] = useState(true);
    // Estado para manejar los errores
    const [error, setError] = useState(null);

    // useEffect para realizar la llamada a la API cuando el componente se monta
    useEffect(() => {
        // Define la función asíncrona para buscar los datos
        const fetchData = async () => {
            try {
                // Reinicia el estado antes de la llamada
                setLoading(true);
                setError(null);

                const response = await fetch(apiUrl);

                // Verifica si la respuesta de la red fue exitosa
                if (!response.ok) {
                    throw new Error(`Error HTTP: ${response.status} ${response.statusText}`);
                }

                const jsonData = await response.json();

                // Asegúrate de que la respuesta sea un array
                if (Array.isArray(jsonData)) {
                    setData(jsonData);
                } else {
                    console.warn("La respuesta de la API no fue un array:", jsonData);
                    setData([]); // Establece un array vacío para evitar errores posteriores
                    // Opcionalmente, podrías lanzar un error aquí si siempre esperas un array
                    // throw new Error("Formato de datos inesperado.");
                }

            } catch (err) {
                // Captura cualquier error durante el fetch o procesamiento
                console.error("Error al obtener los datos:", err);
                setError(err);
                setData([]); // Limpia los datos en caso de error
            } finally {
                // Se ejecuta siempre, haya habido éxito o error
                setLoading(false); // Marca que la carga ha terminado
            }
        };

        fetchData(); // Llama a la función para obtener los datos

        // El array vacío [] como segundo argumento de useEffect significa
        // que este efecto se ejecutará solo una vez, similar a componentDidMount
    }, []); // El array de dependencias vacío es crucial aquí

    // --- Renderizado condicional ---

    // Muestra un mensaje mientras carga
    if (loading) {
        return <p>Cargando datos...</p>;
    }

    // Muestra un mensaje si hubo un error
    if (error) {
        return <p>Error al cargar los datos: {error.message}</p>;
    }

    // Muestra un mensaje si no hay datos (después de cargar y sin errores)
    if (data.length === 0) {
        return <p>No se encontraron datos.</p>;
    }

    // --- Renderizado de la tabla (si hay datos) ---

    // Obtiene las cabeceras de la tabla a partir de las claves del primer objeto
    // Asegúrate de que 'data' tenga al menos un elemento antes de hacer esto
    const headers = Object.keys(data[0]);

    return (
        <div>
            <h1>Tabla de Items (React)</h1>
            <table>
                <thead>
                    <tr>
                        {/* Mapea los nombres de las cabeceras a elementos <th> */}
                        {headers.map(header => (
                            <th key={header}>{header}</th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {/* Mapea cada objeto en el array de datos a una fila <tr> */}
                    {data.map((item, index) => (
                        // Es MUY recomendable usar un ID único del item como 'key' si está disponible
                        // ej: <tr key={item.id}>. Usar el index es un fallback.
                        <tr key={item.id || index}>
                            {/* Para cada fila, mapea las cabeceras para obtener el valor de la celda */}
                            {headers.map(header => (
                                <td key={`${item.id || index}-${header}`}>
                                    {/* Muestra el valor de la celda. Convierte a String por si acaso */}
                                    {/* Maneja valores null/undefined para que no den error */}
                                    {item[header] !== null && item[header] !== undefined ? String(item[header]) : ''}
                                </td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default TablaDatos; // Exporta el componente para usarlo en otros lugares
