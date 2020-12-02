FROM --platform=$TARGETPLATFORM python:3.9
ADD requirements.txt /code/requirements.txt
ADD pip.conf /etc/pip.conf
WORKDIR /code
RUN pip install --upgrade pip
RUN pip install -r requirements.txt
ADD app.py /code/app.py
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "curl","http://localhost:5000/ping" ]
CMD [ "python" ,"app.py"]