# Yandex Container Registry Lifecycle Tool
Инструмент обхода списка container registry и назначения [политик жизненного цикла](https://cloud.yandex.ru/docs/container-registry/operations/lifecycle-policy/lifecycle-policy-create) в промышленных масштабах

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
