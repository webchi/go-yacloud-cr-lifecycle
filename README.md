# Серсис который устанавливает политики управления Docker Regitry 


## Список необходимых переменных окружения 

````bash
YANDEX_FOLDER_ID
YANDEX_QAUTH_TOCKEN
````

## Проверка работы приложения 
````bash

yc container repository list
yc container repository lifecycle-policy list --repository-name crpbcv0kq3k8f2813aha/admin

sergeykletsov@MacBook-Pro-Sergej-2 kubespray_terraform_yandex_cloud % yc container repository lifecycle-policy list --repository-name crpbcv0kq3k8f2813aha/admin
+----------------------+------+----------------------+--------+---------------------+--------------+
|          ID          | NAME |    REPOSITORY ID     | STATUS |       CREATED       | DESCRIPTION  |
+----------------------+------+----------------------+--------+---------------------+--------------+
| crp80smlailff7p4rfja | test | crpqtebedva3i710qm5j | ACTIVE | 2022-04-14 18:22:26 | for testing  |
+----------------------+------+----------------------+--------+---------------------+--------------+

````

