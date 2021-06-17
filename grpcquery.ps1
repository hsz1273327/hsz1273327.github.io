grpcurl -d="{\"model_spec\":{\"name\":\"half_plus_two\",\"version\":123},\"inputs\":{\"x\":{
    'dtype':'DT_FLOAT',
    'tensor_shape':{
    'dim':[
    {
    'size':'3',
    },
    {
    'size':'1',
    }
    ],
    float_val:[1.0, 2.0, 5.0]
    },
}},\"output_filter\":[\"y\"]}" -plaintext -insecure  localhost:8500 PredictionService/Predict