import grpc
import tensorflow as tf
from tensorflow_serving.apis import predict_pb2, prediction_service_pb2_grpc


channel = grpc.insecure_channel("localhost:8500")

# Create the PredictionServiceStub
stub = prediction_service_pb2_grpc.PredictionServiceStub(channel)

# Create the PredictRequest and set its values
req = predict_pb2.PredictRequest()
req.model_spec.name = 'half_plus_two'
req.model_spec.signature_name = ''

# Convert to Tensor Proto and send the request
# Note that shape is in NHWC (num_samples x height x width x channels) format
tensor = tf.make_tensor_proto()
req.inputs["x"].CopyFrom(tensor)  # Available at /metadata

# Send request
response = stub.Predict(req)

# Handle request's response
output_tensor_proto = response.outputs["dense_2"]  # Available at /metadata
