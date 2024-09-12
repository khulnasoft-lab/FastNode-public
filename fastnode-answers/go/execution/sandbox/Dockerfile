FROM python:3.7.3-stretch

WORKDIR /fastnode/app

COPY packages.txt .
RUN pip install --no-cache-dir -r packages.txt

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY sample-files ./sample-files

COPY entrypoint.py .

WORKDIR /fastnode/run
ENTRYPOINT ["python", "/fastnode/app/entrypoint.py"]
