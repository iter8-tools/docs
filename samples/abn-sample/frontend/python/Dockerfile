
FROM python:3.9-alpine
WORKDIR /
COPY samples/abn-sample/frontend/python .
RUN pip install -r requirements.txt

ENV FLASK_APP main.py
ENTRYPOINT ["flask", "run", "--port", "8090"]