FROM ubuntu:22.04

RUN apt update && apt upgrade -y
RUN apt install nodejs npm -y
RUN apt install default-jdk default-jre -y
RUN mkdir /swagger-codegen
WORKDIR /swagger-codegen
RUN wget https://repo1.maven.org/maven2/io/swagger/codegen/v3/swagger-codegen-cli/3.0.20/swagger-codegen-cli-3.0.20.jar -Oswagger-codegen.jar
RUN apt install python3 python3-pip -y
RUN yes | pip3 install pyyaml --break-system-packages
COPY ./openapi.yaml .
RUN echo >script.py "import sys, yaml, json\nwith open('openapi.yaml','r',encoding='utf-8') as f:\n    yaml_data=yaml.safe_load(f)\
    \nwith open('openapi.json','w',encoding='utf-8') as f:\n    json.dump(yaml_data,f,indent=2)\n"
RUN python3 script.py
RUN mkdir -p /opt/www
RUN touch /opt/www/.swagger-codegen-ignore
RUN java -jar swagger-codegen.jar generate \
    -i./openapi.json \
    -l nodejs-server \
    -o /opt/www
WORKDIR /opt/www
RUN npm i
CMD ["npm", "run", "start"]
EXPOSE 8080
