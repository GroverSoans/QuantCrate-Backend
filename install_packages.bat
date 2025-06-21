@echo off
echo Installing packages for BTC-LSTM notebook...
"C:\Users\Grover Soans\AppData\Local\Programs\Python\Python310\python.exe" -m pip install --upgrade pip
"C:\Users\Grover Soans\AppData\Local\Programs\Python\Python310\python.exe" -m pip install -r services\ml-prediction-service\requirements.txt
echo Package installation completed!
pause 