# Khti timetable bot
[![coverage report](https://gitlab.com/qulaz/khti_timetable_bot/badges/master/coverage.svg?job=test_bot)](https://gitlab.com/qulaz/khti_timetable_bot/-/commits/master/bot)

Чат-бот Хакасского технического института на платформе ВКонтакте, который помогает 
студентам следить за расписанием.

## Что умеет бот?

* [x]   Показывать расписание:
    * [x]   На сегодняшний/завтрашний день
    * [x]   На первую/вторую неделю целиком
    * [x]   На конкретный день недели выбранной недели
* [x]   Определять следующую пару пользователя
* [x]   Определять время до звонка с/на пару, учитывая в какой момент времени пришел запрос
* [x]   Определять номер текущей недели (первая/вторая)
* [x]   Оповещать студентов об изменении расписания:
    * [ ]   Показывать что конкретно поменялось в расписании

## Демонстрация работы бота

<div align="center">
  <img src="showcase.gif">
</div>


## Подготовка сервера к деплою

* Установить `Docker` и `docker-compose`:
    * [Установка Docker](https://docs.docker.com/engine/install/ubuntu/)
    * [Запуск без sudo](https://docs.docker.com/engine/install/linux-postinstall/)
    * [Настройка логов](https://docs.docker.com/config/containers/logging/configure/)
    * [Установка docker-compose](https://docs.docker.com/compose/install/)
* Установить `Gitlab Runner`:
    ```shell script
    sudo curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64
    sudo chmod +x /usr/local/bin/gitlab-runner
    sudo gitlab-runner install --user=ВСТАВИТЬ_ИМЯ_ТЕКУЩЕГО_ПОЛЬЗОВАТЕЛЯ --working-directory=ВСТАВИТЬ_ПУТЬ_ДО_ПАПКИ_ТЕКУЩЕГО_ПОЛЬЗОВАТЕЛЯ
    ```
* [Загеристрировать `Gitlab Runner` для проекта с тегом `server`](https://docs.gitlab.com/runner/register/index.html)
