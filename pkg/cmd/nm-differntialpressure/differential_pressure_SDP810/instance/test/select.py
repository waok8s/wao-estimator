#!/usr/bin/python
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base

SQLALCHEMY_DATABASE_URI = 'sqlite:///SDP810.db'
db = SQLAlchemy() 
ma = Marchmallow()

Base = declarative_base()
Base.metadata.reflect(bind=SDP810)

class some_table(Base): 
    __table__ = Base.metadata.tables['sensors']

session = Session()

for row in session.query(sensors.sensorid).filter_by(sensorid="101037B"):
    print(row)
