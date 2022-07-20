# Серсис который устанавливает политики управления Docker Regitry 


## Список необходимых переменных окружения 

````bash
YANDEX_FOLDER_ID
YANDEX_OAUTH_TOKEN
````

## Проверка работы приложения 
````bash

yc container repository list
yc container repository lifecycle-policy list --repository-name 

````
