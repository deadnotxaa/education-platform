# Stage 1. Preparations and design

*Тема: practicum.yandex.ru*

## Функциональные требования

### I. Курсы

**1.1. Информация о курсе**  
Система должна хранить следующую информацию о курсах:
- Название 
- Описание
- Специализация
- Продолжительность
- Цена курса
- Список тем курса
- Информация о преподавателях курса
- Уровень сложности курса (с нуля, продвинутый и т.д.)

**1.2. Информация о темах курса**  
Для каждого курса должен храниться список тем курса. Система должна хранить следующую информацию для каждой темы:
- Название темы
- Описание
- Навыки и технологии осваиваемые в рамках данной темы
- Трудоёмкость темы в часах
- Количетсво промежуточных проектов, входящих в тему

**1.3. Статистика курса**  
- Сколько человек приобрело курс
- Сколько человек были успешно трудоустроены по итогам курса и карьерного трека
- Средняя оценка курса выпускниками

**1.4. Связь с другими объектами**  
- У курса есть состав преподавателей, менторов и ревьюеров


### II. Сотрудники

**2.1. Информация о сотрудниках**  
Система должна иметь возможность регистрировать сотрудников платформы:
- Преподаватель
- Ментор
- Ревьюер
- Поддержка
- SMM-менеджер

Для каждого необходимо указать:
- Имя и фамилия
- Дата рождения
- Должность


### III. Расписание курсов

**3.1. Календарь курсов**  
Система должна хранить информацию о расписании проведения курсов:
- Дата начала курса
- Дата дедлайна для регистрации на прохождение (покупки)
- Информация о наличии мест

**3.2. Связь с другими объектами**  
- Каждая отметка в календаре курсов должна соответствовать непосредственно существующему курсу


### IV. Преподаватели

**4.1. Информация о преподавателе**  
- Имя и фамилия
- Текущее место работы
- Количество лет опыта в сфере в общем; конкретно в этой специализации
- Фотография

Например:

***Иван Иванов***  
*Старший разработчик в Тындекс <любой сервис>. Программирует 10 лет, последние 5 — на Brainfuck. Разрабатывает архитектуру очень медленных систем. Входит в Brainfuck‑комитет Тындекса.*

**4.2. Принадлежность к курсу**  
- Должна быть информация, к каким курсам относится преподаватель


### V. Взаимодействие с пользователями

**5.1. Регистрация пользователей**  
Система должна позволять новым пользователям регистрироваться. При этом пользователь указывает:
- Имя и фамилия
- Дата рождения
- Электронная почта
- Пароль

**5.2. Аутентификация пользователей**  
Зарегистрированные пользователи должны иметь возможность войти в систему, используя электронную почту и пароль.

**5.3. Управление профилем**  
Пользователь может редактировать свой профиль; добавлять информацию о себе:
- Фотография
- Номер телефона
- СНИЛС


### VI. Покупка курсов

**6.1. Выбор курса**  
Пользователь должен иметь возможность купить любой доступный курс. Система должна:
- Отображать актуальные и доступные курсы для покупки
- Позволять пользователю видеть информацию об актуальной цене курса

**6.2. Виды курсов**  
Должна быть возможность приобрети курс:
- Базовый
- С расширенным карьерным треком
- Льготный

**6.3. Информация о заказе**  
Система должна сохранять информацию по каждому купленному курсу (в одном заказе может быть только один курс):
- ID пользователя
- Дата и время заказа
- Вид курса
- Итоговая сумма заказа
- Статус (оплачен/неоплачен)


### VII. Карьерный центр

**7.1. Информация о пользователях, находящихся в ведении карьерного центра**  
Система должна хранить следующую информацию о пользователях, находящихся в её ведении:
- Имя и фамилия
- Пройденный курс
- Актуальное резюме
- Дата начала карьерного сопровождения
- Период сопровождения (расширенный или нет)

**7.2. Информация о компаниях-партнёрах**  
Система должна хранить информацию о компаниях партнёрах, с которыми заключены соглашения
- Название компании (короткое и полное юридическое наименование)
- Количество выпускников, принятых в компанию партнёр
- Требования к кандидатам
- Статус соглашения (активное или нет)


### VIII. Блог

**8.1. Информация о статьях**  
Система должна хранить информацию обо всех опубликованных статьях, каждая из которых включает:
- Название
- Автор статьи
- Дата публикации
- Тема
- Примерное время на прочтение статьи
- Обложка
- Содержание

**8.2. Связь с другими объектами**  
Автор статьи является конкретным SMM-менеджером, зарегистрированным в системе

## ERD

```plantuml
@startuml

' Courses
entity course {
    *course_id : serial <<PK>>
    --
    name : varchar
    description : text
    specialization_id : integer <<FK>>
    duration : integer
    price : integer
    difficulty_level_id : integer <<FK>>
}

entity difficulty_level {
    *id : serial <<PK>>
    --
    name : varchar
    description : text
}

entity course_specialization {
    *id : serial <<PK>>
    --
    name : varchar
    description : text
}

entity course_topic {
    *id : serial <<PK>>
    --
    name : varchar
    description : text
    technologies : text
    labor_intensity_hours : smallinteger
    projects_number : smallinteger
}

entity course_topic_association {
    *course_id : integer <<FK>>
    *topic_id : integer <<FK>>
}

entity project {
    *project_id : serial <<PK>>
    --
    topic_id : integer <<FK>>
    name : varchar
    description : text
}

' Users
entity user {
    *account_id : serial <<PK>>
    --
    name : varchar
    surname : varchar
    birthdate : date
    email : varchar
    hashed_password : varchar
    profile_picture_url : varchar [nullable]
    phone_number : varchar [nullable]
    snils_number : varchar [nullable]
}

' Employees
entity role {
    *id : serial <<PK>>
    --
    name : varchar
}

entity employee {
    *id : serial <<PK>>
    --
    user_id : integer <<FK>>
    role_id : integer <<FK>>
}

' Schedule
entity course_calendar {
    *id : serial <<PK>>
    --
    course_id : integer <<FK>>
    start_date : date
    end_sales_date : date
    remaining_places : smallint
}

' Teachers
entity teacher {
    *employee_id : integer <<PK>> <<FK>>
    --
    work_place : varchar
    overall_experience : integer
    specialization_experience : integer
}

entity course_teacher {
    *teacher_id : integer <<FK>>
    *course_id : integer <<FK>>
}

' Course Purchasing
entity purchase {
    *purchase_id : serial <<PK>>
    --
    user_id : integer <<FK>>
    course_id : integer <<FK>>
    purchase_date : datetime
    course_type_id : smallinteger <<FK>>
    total_price : integer
    purchase_status : payment_status
}

entity course_type {
    *id : smallserial <<PK>>
    --
    type_name : varchar
    discount : smallinteger
}

' Career Centre
entity career_center_student {
    *id : serial <<PK>>
    --
    user_id : integer <<FK>>
    course_id : integer <<FK>>
    cv_url : varchar
    career_support_start : date
    support_period : integer
}

entity partner_company {
    *company_id : serial <<PK>>
    --
    short_name : varchar
    full_name : varchar
    hired_graduates_count : integer
    requirements : text
    agreement_status : boolean
}

entity job_application {
    *id : serial <<PK>>
    --
    student_id : integer <<FK>>
    company_id : integer <<FK>>
    application_date : date
    status : varchar
}

' Blog
entity blog_post {
    *post_id : serial <<PK>>
    --
    author_id : integer <<FK>>
    title : varchar
    publication_date : datetime
    topic : varchar
    reading_time_minutes : integer
    cover_image_url : varchar
    content : text
}

entity tag {
    *tag_id : serial <<PK>>
    --
    name : varchar
}

entity post_tag {
    *post_id : integer <<FK>>
    *tag_id : integer <<FK>>
}

' Reviews
entity course_review {
    *review_id : serial <<PK>>
    --
    course_id : integer <<FK>>
    user_id : integer <<FK>>
    rating : smallinteger
    comment : text
    review_date : datetime
}

' Certificates
entity certificate {
    *certificate_id : serial <<PK>>
    --
    user_id : integer <<FK>>
    course_id : integer <<FK>>
    issue_date : date
}

' Relations
course_specialization::id ||--o{ course::specialization_id
difficulty_level::id ||--o{ course::difficulty_level_id
course::course_id ||--o{ course_topic_association::course_id
course_topic::id ||--o{ course_topic_association::topic_id
course_topic::id ||--o{ project::topic_id
user::account_id ||--o{ employee::user_id
role::id ||--o{ employee::role_id
course::course_id ||--o{ course_calendar::course_id
employee::id ||--o{ teacher::employee_id
teacher::employee_id ||--o{ course_teacher::teacher_id
course::course_id ||--o{ course_teacher::course_id
course_type::id ||--o{ purchase::course_type_id
user::account_id ||--o{ purchase::user_id
course::course_id ||--o{ purchase::course_id
user::account_id ||--o{ career_center_student::user_id
course::course_id ||--o{ career_center_student::course_id
career_center_student::id ||--o{ job_application::student_id
partner_company::company_id ||--o{ job_application::company_id
employee::id ||--o{ blog_post::author_id
blog_post::post_id ||--o{ post_tag::post_id
tag::tag_id ||--o{ post_tag::tag_id
course::course_id ||--o{ course_review::course_id
user::account_id ||--o{ course_review::user_id
user::account_id ||--o{ certificate::user_id
course::course_id ||--o{ certificate::course_id

@enduml
```
