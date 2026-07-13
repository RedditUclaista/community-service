# Community Microservice

Este es el microservicio encargado de gestionar las comunidades (equivalente a los subreddits).

## ¿Qué hace?
- Maneja la creación, configuración y listado de comunidades.
- Gestiona la suscripción y roles de los usuarios dentro de comunidades específicas.

## ¿Cómo levantarlo?
Este servicio utiliza un Dockerfile estándar de Go y depende de la infraestructura compartida del API Gateway.

1. Asegúrate de tener levantado el Gateway (y su red de Docker) primero.
2. Desde este directorio (`backend/community`), ejecuta:
   ```bash
   docker-compose up -d
   ```
3. El contenedor se compilará, se unirá a la red `lab3_shared_network` y operará en su puerto (`10001`).
