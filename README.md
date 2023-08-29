<h1>Сервис динамического сегментирования пользователей</h1>
<p>Этот проект представляет собой сервис для динамического сегментирования пользователей, разработанный на языке Golang. Сервис позволяет создавать, изменять, удалять сегменты и добавлять или удалять пользователей из сегментов. Задачей сервиса является централизация работы с проводимыми экспериментами и предоставление API для управления сегментами и пользователями.</p>
<h2>Основные требования</h2>
<h3>Технические требования:</h3>
<ul>
  <li>Язык разработки: Golang.</li>
  <li>Использование фреймворков и библиотек на ваше усмотрение.</li>
  <li>Реляционная СУБД: MySQL или PostgreSQL.</li>
  <li>Docker и docker-compose для поднятия и развертывания dev-среды.</li>
  <li>HTTP API с форматом JSON для обмена данными.</li>
</ul>
<h3>Установка и запуск:</h3>
<ol>
  <li>Склонируйте репозиторий:</li>
</ol>
<pre><code>git clone github.com/Dani1Kruglov/avito-test-task</code></pre>
<ol start="2">
  <li>Запустите контейнеры с БД и сервисом:</li>
</ol>
<pre><code>docker-compose up -d</code></pre>
<ol start="3">
  <li>Запустите сервис:</li>
</ol>
<ol start="4">
  <li>Переименуйте example-config.hcl в config.hcl и заполните свои данные бд</li>
</ol>
<pre><code>go run main.go</code></pre>
<h3>Использование API:</h3>
<h4>Создание сегмента (метод: StoreSegment с переданным полем: название сегмента в JSON формате)</h4>
<h4>Удаление сегмента (метод: DeleteSegment с переданным полем: название сегмента в JSON формате)</h4>
<h4>Добавление пользователя в сегмент (метод: AddUserToSegments с переданными полями: id пользователя и названия сегментов в JSON формате)</h4>
<h4>Удаление пользователя из сегмента (метод: DeleteUserSegments с переданными полями: id пользователя и названия сегментов в JSON формате)</h4>
<h4>Получение активных сегментов пользователя (метод: GetUserSegments с переданным полем: id пользователя в JSON формате)</h4>
<h4>Также, в качестве дополнения, созданы методы CRUD для работы с пользователями (методы находятся в internal/storage/users.go). Передача в методы для пользователей происходит также в JSON формате.</h4>
<h2>Опциональные задания</h2>
<h3>Доп. задание 1: История попадания в сегмент</h3>
<ul>
  <li>Реализовано сохранение истории попадания/выбывания пользователя из сегмента.</li>
  <li>Для получения отчета за определенный период (от нескольких секунд) обратитесь к методу GetUserSegmentsInCSVFile
</ul>
<h3>Доп. задание 2: TTL для пользователей</h3>
<ul>
  <li>Реализована возможность задавать TTL (время автоматического удаления пользователя из сегмента) при добавлении пользователя в сегмент.</li>
  <li>Для добавления пользователя в сегмент на некоторое время используйте метод AddUserToSegmentsForWhile. В метод необходимо передать данные пользователя и сегменты в JSON формате, а также продолжительность хранения в бд.</li>
</ul>
<h3>Доп. задание 3: Автоматическое добавление пользователей</h3>
<ul>
  <li>Добавлена опция указания процента пользователей, которые будут попадать в сегмент автоматически при создании сегмента.</li>
  <li>При создании сегмента с процентом пользователей используйте метод AddSomeUsersToSegment. В метод необходимо передать данные сегментов в JSON формате, а также процент от общего количества пользователей,которые будут иметь доступ к этим сегментам .</li>
</ul>
<p>Этот README содержит описание основных функций, инструкции по установке и использованию, а также информацию о выполненных опциональных заданиях.</p>
