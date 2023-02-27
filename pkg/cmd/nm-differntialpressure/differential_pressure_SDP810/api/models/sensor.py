from api.database import db, ma
import json

class Sensor(db.Model):
  __tablename__ = 'sensors'

  sensorid = db.Column(db.String(10), primary_key=True, unique=True, nullable=False)
  pressure = db.Column(db.Float, nullable=True)
  temperature = db.Column(db.Float, nullable=True)

  def __repr__(self):
    return '<Sensor %r>' % self.name

  def getSensorList():

    # select * from sensors
    sensor_list = db.session.query(Sensor).all()
    if sensor_list == None:
      return []
    else:
      return sensor_list

  def registSensor(sensor):
    record = Sensor(
      sensorid = sensor['sensorid'],
      pressure = sensor['pressure'],
      temperature = sensor['temperature']
    )

    # insert into sensors(sensorid, presssure, temperature) values(...)
    db.session.add(record)
    db.session.commit()

    return sensor 


  def getSensor(sensorid):
    # select * from sensors where sensorid = sensorid
    sensor = db.session.query(Sensor).filter(Sensor.sensorid == sensorid)

    if sensor == None:
        return []
    else:
        return sensor

class SensorSchema(ma.SQLAlchemyAutoSchema):
    class Meta:
      model = Sensor 
      fields = ('sensorid', 'pressure', 'temperature')

