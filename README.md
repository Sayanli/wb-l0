# wb-l0

Все основные требования были выполнены.

PosgreSQL была развернута в docker контейнере.

Реализована подписка на канал nats-streaming. Запись идёт и в кэш, и в базу данных postgreSQL. Был сделан небольшой скрипт, который позволяет отправлять json в канал.

При падении приложения кэш восстанавливается из базы данных.

Реализован http сервер с использованием fiber.

Для удобства был сделан простой интерфейс, который позволяет искать заказ по UID.
