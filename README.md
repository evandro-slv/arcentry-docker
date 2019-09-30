# arcentry-docker

Tool used to monitor running docker containers and report cpu and memory usage into Arcentry.

# Usage

First, download the source code and compile it:

    go build -o arcentry-docker
   
This will generate a `arcentry-docker` executable file in the directory, you can then copy it to the desired folder:

    mv arcentry-docker /usr/bin/

There is a required config file needed to run the app, the file should be located at `$HOME/.arcentry-docker.yaml` or you can send the path as an argument:

    arcentry-docker -c config.yml


Config template:    
    
    config:
      arcentry:
        apiKey: {API_KEY}
        docId: {DOC_ID}
      watch:
        interval: "30s"
      containers:
        - id:
            chart:
              cpu: {OBJ_ID}
              memory: {OBJ_ID}
            text:
              cpu: {OBJ_ID}
              memory: {OBJ_ID}

- `API_KEY`: the arcentry API key
- `DOC_ID`: your document id
- `OBJ_ID`: the object id to be updated

The objects inside the `chart` element will display a char with the last 10 cpu and memory reads, while the `text` element will send only the value (10.5% for example)

`watch`: is the interval it will wait to send the data

Example:

    config:
      arcentry:
        apiKey: "37ad4c37179b4233ad3472ce1fedd9c537ad4c37179b4233ad3472ce1fedd9c5"
        docId: "37ad4c37-179b-4233-ad34-72ce1fedd9c5"
      watch:
        interval: "30s"
      containers:
        - b7252fa2:
            chart:
              cpu: b7252fa2-0710-4af5-aae3-75c5a2e8cf9c 
              memory: 9246067b-edca-405f-8271-47628eeaaf4a 
            text:
              cpu: c7fe6a6e-e3cc-4ec6-92e6-6201a9844afa
              memory: 37cf0f04-bbb9-41e3-b7f3-6254191eb917
        - 640c279e8151:
          chart:
            cpu: b7252fa2-0710-4af5-aae3-75c5a2e8cf9c 
            memory: 9246067b-edca-405f-8271-47628eeaaf4a 
          text:
            cpu: c7fe6a6e-e3cc-4ec6-92e6-6201a9844afa
            memory: 37cf0f04-bbb9-41e3-b7f3-6254191eb917

The project is in development phase, any help is appreciated.