# Простой пример работы прометеуса
 Сначала описываем docker-compose, указываем сервисы, в нашем случае образ прометеуса.

После этого переходим в файлы prometheus.yml и описываем там конфиг

Также была установлена утилита ghz для нагрузочного тестирования  

Одна тонкость по запуску прометеуса, мы указали targets: ["host.docker.internal:2112"], потому что если указать просто localhost, то собрать метрики не получится.
Потому что прометеус сам запускается в докере, и мы ему прямо указываем, что нужен этот порт.

# Теперь определяем алерты
Переходим в alerts.yml и описываем алерты.
И теперь если условие, которые мы задали не выполнится, то придет алерт.

# Добавляем теперь Response metrics, чтобы можно посмотреть сколько запросов завершились успешно и с ошибками.
Теперь мы можем увидеть в прометеусе эти ответы.
Но для того, чтобы увидеть еще лэйбл с ошибкой, у нас есть команда в makefile - grpc-error-load-test
И таким образом увидим, что 

# Добавляем Grafana для отображения метрик
Теперь мы добавим дашборд, на котором будут отображаться наши метрики, для более наглядного отображения.
Для этого, мы запустим контейнер с Grafana, далее перейдем по порту 3000 и зайдем в саму Grafana.
После этого переходим в Data Sources и добавляем прометеус. 
Указываем имя Prometheus и http-port, здесь есть нюанс, так как у нас и прометеус и графана запущены в докере,
Поэтому в рамках докер сети, мы обращаемся к прометеусу по адресу http://prometheus:9090
Теперь мы наконце можем создать наш дашборд.

Добавил файл grafana_dashboard.json, который можно импортировать в графану, для отображения настроенных дашбордов.