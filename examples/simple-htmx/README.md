## Simple HTMX demo

This demo demonstrates embedding of posts using Samara's HTMX support.
Additionally, this also includes OTEL configuration so that you can see what is happening in the background.

```
# Start Alloy and a frontend server
docker compose up

# Start backend
cd $ROOT_FOLDER
just run-demo
```

Now open your browser and go to http://localhost:8888.
You will see a post being loaded and Alloy receiving some traces in the terminal.
