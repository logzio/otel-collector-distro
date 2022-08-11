# Test PerformanceResults
Started: Mon, 04 Jul 2022 11:29:15 +0300

Test                                    |Result|Duration|CPU Avg%|CPU Max%|RAM Avg MiB|RAM Max MiB|Sent Items|Received Items|
----------------------------------------|------|-------:|-------:|-------:|----------:|----------:|---------:|-------------:|
IdleMode                                |PASS  |     16s|     1.1|     5.2|         25|         36|         0|             0|
Log10kDPS/OTLP                          |PASS  |     15s|    21.6|    23.1|         43|         63|    149900|        149900|
Log10kDPS/OTLP-HTTP                     |PASS  |     15s|    17.5|    18.1|         37|         53|    149900|        149900|
Log10kDPS/filelog                       |PASS  |     15s|    23.6|    24.3|         43|         61|    150000|        150000|
Log10kDPS/kubernetes_containers         |PASS  |     15s|    47.8|    48.4|         46|         66|    150000|        150000|
Log10kDPS/k8s_CRI-Containerd            |PASS  |     15s|    45.1|    45.9|         45|         65|    150000|        150000|
Log10kDPS/k8s_CRI-Containerd_no_attr_ops|PASS  |     15s|    36.6|    38.0|         45|         65|    150000|        150000|
Log10kDPS/CRI-Containerd                |PASS  |     15s|    23.7|    24.9|         45|         64|    150000|        150000|
Log10kDPS/syslog-tcp-batch-1            |PASS  |     15s|    27.8|    29.0|         39|         56|    149900|        149900|
Log10kDPS/syslog-tcp-batch-100          |PASS  |     15s|    18.5|    19.4|         39|         56|    149900|        149900|
Log10kDPS/tcp-batch-1                   |PASS  |     15s|    27.8|    28.6|         41|         58|    149900|        149900|
Log10kDPS/tcp-batch-100                 |PASS  |     15s|    17.9|    19.0|         39|         56|    149900|        149900|
Trace10kSPS/OTLP-Logzio                 |PASS  |     15s|    30.1|    31.2|         48|         68|    149700|        149700|
Trace10kSPS/OTLP-HTTP-Logzio            |PASS  |     15s|    27.3|    28.9|         41|         59|    149900|        149900|
Trace10kSPS/Jaeger-Logzio               |PASS  |     15s|    34.7|    36.2|         42|         59|    149900|        149900|
Trace10kSPS/Zipkin-Logzio               |PASS  |     15s|    31.9|    32.9|         37|         54|    149900|        149900|
Trace10kSPS/OTLP-HTTP                   |PASS  |     15s|    17.2|    18.1|         35|         51|    149900|        149900|
TraceNoBackend10kSPS/NoMemoryLimit      |FAIL  |     25s|    34.7|    43.3|         90|        146|    149960|             0|signal: killed
TraceNoBackend10kSPS/MemoryLimit        |PASS  |     23s|    44.0|    53.5|         42|         59|    148980|             0|
Trace1kSPSWithAttrs/0*0bytes            |PASS  |     15s|    26.1|    28.8|         41|         62|     15000|         15000|
Trace1kSPSWithAttrs/100*50bytes         |PASS  |     15s|    50.6|    51.8|         44|         63|     14990|         14990|
Trace1kSPSWithAttrs/10*1000bytes        |PASS  |     15s|    47.0|    47.7|         44|         63|     14990|         14990|
Trace1kSPSWithAttrs/20*5000bytes        |PASS  |     15s|   120.7|   121.7|         54|         77|      5450|          5450|
TraceBallast1kSPSWithAttrs/0*0bytes     |PASS  |     15s|    25.1|    26.2|         74|        125|     14990|         14990|
TraceBallast1kSPSWithAttrs/100*50bytes  |PASS  |     15s|    40.8|    42.6|        575|       1064|     15000|         15000|
TraceBallast1kSPSWithAttrs/10*1000bytes |PASS  |     15s|    35.6|    36.1|        378|        839|     14990|         14990|
TraceBallast1kSPSWithAttrs/20*5000bytes |PASS  |     15s|   104.1|   104.8|        737|       1099|      3920|          3920|
TraceBallast1kSPSAddAttrs/0*0bytes      |PASS  |     15s|    25.1|    28.0|         78|        136|     15000|         15000|
TraceBallast1kSPSAddAttrs/100*50bytes   |PASS  |     15s|    41.4|    43.9|        593|       1061|     14990|         14990|
TraceBallast1kSPSAddAttrs/10*1000bytes  |PASS  |     15s|    37.6|    40.1|        502|       1043|     14990|         14990|
TraceBallast1kSPSAddAttrs/20*5000bytes  |PASS  |     15s|   104.4|   106.2|        711|       1094|      4560|          4560|
TraceAttributesProcessor/OTLP           |PASS  |      1s|     0.0|     0.0|          0|          0|         3|             3|
TraceAttributesProcessor/OTLP-Logzio    |FAIL  |      0s|     0.0|     0.0|          0|          0|         1|             1|

Total duration: 493s
