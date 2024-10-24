# FROM ubuntu:latest

# RUN apt update && apt upgrade -y
# RUN apt install -y wget nodejs npm default-jdk default-jre python3 python3-pip

# RUN mkdir /swagger-codegen
# WORKDIR /swagger-codegen

# RUN wget https://repo1.maven.org/maven2/io/swagger/codegen/v3/swagger-codegen-cli/3.0.20/swagger-codegen-cli-3.0.20.jar -O swagger-codegen.jar

# RUN pip3 install pyyaml --break-system-packages

# COPY ./openapi.yaml .

# RUN echo "import yaml, json\nwith open('openapi.yaml', 'r', encoding='utf-8') as f:\n    yaml_data = yaml.safe_load(f)\nwith open('openapi.json', 'w', encoding='utf-8') as f:\n    json.dump(yaml_data, f, indent=2)" > script.py
# RUN python3 script.py

# RUN mkdir -p /opt/www
# RUN touch /opt/www/.swagger-codegen-ignore
# RUN java -jar swagger-codegen.jar generate \
#     -i ./openapi.json \
#     -l nodejs-server \
#     -o /opt/www

# WORKDIR /opt/www
# RUN npm install oas3-tools@2.0.2
# RUN npm install express-validator@5.3.1 --save-exact
# RUN npm install

# CMD ["npm", "run", "start"]
# EXPOSE 8080

FROM nginx:alpine

# Создаем директорию для хранения Swagger UI и копируем в нее спецификацию OpenAPI
RUN mkdir /swagger_docs
WORKDIR /swagger_docs

# Копируем стандартное содержимое Swagger UI в директорию сервера
RUN wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v4.18.1.tar.gz && \
tar -xzvf v4.18.1.tar.gz && \
cp -r swagger-ui-4.18.1/dist/* . && \
rm -rf v4.18.1.tar.gz swagger-ui-4.18.1

# Настройка Swagger UI для чтения вашего файла openapi.yaml
RUN echo "window.onload = function() { window.ui = SwaggerUIBundle({ url: '/openapi.yaml', dom_id: '#swagger-ui', presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset], layout: 'StandaloneLayout' }); };" >swagger-initializer.js

COPY ./openapi.yaml .
COPY ./nginx.conf .

# RUN ls -la && exit 2

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;", "-c./nginx.conf"]
