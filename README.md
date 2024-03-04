# lorem Ipsum
Lorem Ipsum is a website where you can upload meme and share it for the whole internet to see, it has enough feature. Like upload image, react, comment,
and delete your post. Its pretty much. I think i can improve it, But here it is. 

## Feature
1. View post that user share
2. React and Comment to that post
3. Account authentication
4. Uploading image and delete it

## How to host it yourself
1. First create this compose.yaml
```yaml
services:
  web:
    image: muhhae/lorem-ipsum
    ports:
      - "8080:8080"
    environment:
      - MONGO_URI=mongodb://mongodb:27017/data
      - PORT=8080
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - mongodb
  mongodb:
    image: mongo
    ports:
      - '27017:27017'
    volumes:
      - mongodb-data:/data/db
volumes:
  mongodb-data:
```
2. Run 
```bash 
docker compose up -d
# reload service on change or just update lorem ipsum image
docker compose up -d --pull=muhhae/lorem-ipsum
```

