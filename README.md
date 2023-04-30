# todolist
1. Есть две роли: admin & user
2. Есть пользователи. У них есть ID, name, role, created_at, login, password.
3. Метод добавления пользователя. Пользователь вводит name, login, password. ID задается базой, role задается бизнес логикой и по умолчанию равна user
4. Метод получения пользователя по ID
5. Сделать так, чтобы пароль в базу попадал хэшированным (пакет bcrypt)
6. Убрать поле password при получении пользователя
7. Добавить аутентификацию на cookies
    - Научиться устанавливать cookie в ответ на запрос
    - В обработчике Login принимать запрос методом POST, где в Body будет объект в формате json с полями login и password
    - Ищем в базе пользователя с таким Login и проверяем, совпадают ли пароли.
    - Генерируем id для сессии(с помощью библиотеки google/uuid, почитать, что такое uuid)
    - Сохраняем в базу пару сессию (session id, user id, created_at, expired_at)
    - Сохранить в ответ пользователю cookie с id и временем окончания сессии
8. Добавить логику, что админы могут получать информацию про любого пользователя, а юзеры только свою
9. Добавить метод удаления пользователя по id(пользователь может удалить только себя)
10. Удалять пользователя по id могут только админы
11. Заменить удаление пользователя на soft delete (не удалять запись из БД, а помечать ее как удаленную)
12. Не показывать удаленных пользователей при поиске
13. Перенести проверку сессии в middleware через контекст
14. Добавить метод добавления задачи юзером
15. Добавить метод завершения задачи
16. Добавить метод получения списка активных задач
17. Задача имеет название, ид юзера, статус, ид задачи. 
18. Редактирование и удаление задач через ИД
19. При создании пользователя добавить валидацию (name пользователя только буквы 3-40 символов, латиница, login 3-15 символов, буквы, цифры, точки, нижн. подчеркивания)
20. Добавить функцию валидации пароля. Пароль должен быть не короче 8 символов, латиница, содержать буквы разных регистров и цифры
21. Написать тесты на функции валидации
22. По умолчанию создавать задачу со статусом-константой(active, completed, deleted) и id пользователя
