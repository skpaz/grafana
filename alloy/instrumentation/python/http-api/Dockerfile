FROM python:latest
WORKDIR /app
COPY __init__.py cities.json main.py requirements.txt  ./
RUN python3 -m pip install -r requirements.txt
CMD ["fastapi", "run", "./main.py", "--port", "8080"]
