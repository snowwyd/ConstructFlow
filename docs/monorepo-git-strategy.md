# Инструкция по организации Git-веток для фронтенд и бэкенд разработки

## Обзор подхода

В нашем репозитории используется монорепозиторий с двумя основными компонентами:

- Папка `/frontend` - для фронтенд-кода
- Папка `/backend` - для бэкенд-кода

Основные ветки репозитория:

- Ветка `master` - стабильная версия продукта
- Ветка `front-dev` - разработка фронтенда
- Ветка `back-dev` - разработка бэкенда

## Настройка рабочей среды

### Первоначальная настройка веток

```bash
# Создание ветки для фронтенда
git checkout master
git checkout -b front-dev

# Создание ветки для бэкенда
git checkout master
git checkout -b back-dev
```

Через UI GitHub:

1. Откройте репозиторий на GitHub
2. Перейдите в раздел "Branches"
3. Нажмите "New branch"
4. Введите имя ветки (например, "front-dev")
5. Убедитесь, что выбрана ветка "master" в качестве источника
6. Нажмите "Create branch"
7. Повторите для создания ветки "back-dev"

## Правила для команды разработчиков

### Общие правила

- Ветка `master` защищена, прямые коммиты запрещены
- Все изменения в `master` попадают только через Pull Request
- Каждый Pull Request требует как минимум одного одобрения от другого члена команды

### Правила оформления коммитов

- Формат: `#проект Fixed(что вы сделали). Ref задача`
- В начале указывается название проекта с символом '#'
- Затем глагол, описывающий действие ("Fixed", "Added", "Made", "Implemented" и т.д.)
- В конце добавляется ссылка на задачу с префиксом "Ref"

Примеры:

```
#cf Fixed dependecies. Ref #1
#cf Fixed dependecies. Ref CF-1
#cf Made an authorization server. Ref CF-2
```

### Для фронтенд-разработчиков

- Работайте только в ветке `front-dev`
- Изменяйте только файлы в папке `/frontend`
- Регулярно получайте изменения из `master` для синхронизации с бэкендом

### Для бэкенд-разработчиков

- Работайте только в ветке `back-dev`
- Изменяйте только файлы в папке `/backend`
- Регулярно получайте изменения из `master` для синхронизации с фронтендом

## Рабочий процесс

### Создание feature-веток для новых функций

Для каждой новой функции или задачи следует создавать отдельную ветку (feature branch) из соответствующей ветки разработки:

#### Для фронтенд-задач:

```bash
# Переключиться на ветку фронтенда
git checkout front-dev
git pull origin front-dev

# Создать новую feature-ветку
git checkout -b front-feature-login-form
```

Через UI GitHub:

1. Перейдите на страницу репозитория в GitHub
2. Убедитесь, что выбрана ветка "front-dev" в выпадающем списке веток
3. Нажмите на выпадающий список веток
4. Нажмите "New branch"
5. Введите имя ветки (например, "front-feature-login-form")
6. Убедитесь, что выбрана ветка "front-dev" в качестве источника
7. Нажмите "Create branch"

#### Для бэкенд-задач:

```bash
# Переключиться на ветку бэкенда
git checkout back-dev
git pull origin back-dev

# Создать новую feature-ветку
git checkout -b back-feature-auth-api
```

Через UI GitHub:

1. Перейдите на страницу репозитория в GitHub
2. Убедитесь, что выбрана ветка "back-dev" в выпадающем списке веток
3. Нажмите на выпадающий список веток
4. Нажмите "New branch"
5. Введите имя ветки (например, "back-feature-auth-api")
6. Убедитесь, что выбрана ветка "back-dev" в качестве источника
7. Нажмите "Create branch"

#### Правила именования feature-веток:

- Используйте префикс `front-feature-` или `back-feature-`
- Добавляйте краткое описание функциональности через дефис, например: `login-form`, `user-profile`, `auth-api`
- Для исправления ошибок используйте префиксы `front-fix-` или `back-fix-`

### Разработка новой функциональности

1. Фронтенд и бэкенд-команды обсуждают и согласовывают API/интерфейсы
2. Каждая команда работает в своих feature-ветках

### Работа в feature-ветке (для всех разработчиков)

```bash
# Работа в своей feature-ветке
git checkout front-feature-login-form  # или ваша feature-ветка

# Внести изменения и закоммитить
git add frontend/  # или backend/ для бэкенд-задач
git commit -m "#cf Добавлен компонент авторизации. Ref CF-42"

# Отправить изменения в репозиторий
git push origin front-feature-login-form
```

Через UI GitHub:

1. Убедитесь, что вы работаете в правильной ветке, выбрав её из выпадающего списка веток
2. Для изменения файлов нажмите на файл, который хотите отредактировать
3. Нажмите на иконку карандаша (Edit this file)
4. Внесите необходимые изменения
5. В поле "Commit changes" введите сообщение в формате: `#cf Добавлен компонент авторизации. Ref CF-42`
6. Выберите опцию "Commit directly to the [your-branch-name] branch"
7. Нажмите "Commit changes"

Для добавления новых файлов:

1. Перейдите в нужную папку (например, /frontend)
2. Нажмите кнопку "Add file" → "Create new file" или "Upload files"
3. Создайте или загрузите файл
4. Оформите коммит как описано выше

### Завершение работы в feature-ветке

Когда функция готова:

1. Убедитесь, что код работает и проходит все тесты
2. Создайте Pull Request из вашей feature-ветки в соответствующую ветку разработки (`front-dev` или `back-dev`)

### Слияние feature-ветки в ветку разработки

```bash
# Для фронтенд feature-веток
git checkout front-dev
git pull origin front-dev
git merge front-feature-login-form --no-ff -m "#cf Merge: Добавлен компонент авторизации. Ref CF-42"
git push origin front-dev

# Для бэкенд feature-веток
git checkout back-dev
git pull origin back-dev
git merge back-feature-auth-api --no-ff -m "#cf Merge: Реализован API авторизации. Ref CF-42"
git push origin back-dev
```

Через GitHub UI (предпочтительный метод):

1. Перейдите на страницу репозитория в GitHub
2. Нажмите вкладку "Pull requests"
3. Нажмите кнопку "New pull request"
4. В выпадающем списке "base:" выберите ветку разработки (`front-dev` или `back-dev`)
5. В выпадающем списке "compare:" выберите вашу feature-ветку
6. Нажмите "Create pull request"
7. Заполните заголовок в формате: `#cf Добавлен компонент авторизации. Ref CF-42`
8. Добавьте описание изменений в теле PR
9. Назначьте ревьюеров
10. После получения одобрения (approval) нажмите "Merge pull request"
11. Подтвердите слияние

## Слияние изменений в master

### Через GitHub UI (рекомендуется)

1. Перейдите на страницу репозитория в GitHub
2. Нажмите **Pull requests** → **New pull request**
3. Выберите `base: master` и `compare: front-dev` (или `back-dev`)
4. Нажмите **Create pull request**
5. Заполните заголовок и описание
   - Для фронтенда: "#cf Frontend: [краткое описание]. Ref CF-42"
   - Для бэкенда: "#cf Backend: [краткое описание]. Ref CF-42"
6. Назначьте ревьюеров и дождитесь одобрения
7. Выберите метод слияния **Create a merge commit**
8. Нажмите **Merge pull request**

### Через командную строку (альтернатива)

```bash
# Слияние фронтенд-изменений
git checkout master
git pull origin master
git merge front-dev --no-ff -m "#cf Merge frontend: Добавление функции X. Ref CF-42"
git push origin master

# Слияние бэкенд-изменений
git checkout master
git pull origin master
git merge back-dev --no-ff -m "#cf Merge backend: Реализация API для функции X. Ref CF-42"
git push origin master
```

## Синхронизация рабочих веток после слияния в master

### Для фронтенд-ветки

```bash
git checkout front-dev
git pull origin master
git push origin front-dev
```

Через UI GitHub:

1. Перейдите на страницу репозитория в GitHub
2. Нажмите вкладку "Pull requests"
3. Нажмите кнопку "New pull request"
4. В выпадающем списке "base:" выберите `front-dev`
5. В выпадающем списке "compare:" выберите `master`
6. Нажмите "Create pull request"
7. Заполните заголовок: `#cf Синхронизация front-dev с master. Ref CF-XX`
8. Нажмите "Create pull request"
9. Сразу нажмите "Merge pull request" и подтвердите слияние

### Для бэкенд-ветки

```bash
git checkout back-dev
git pull origin master
git push origin back-dev
```

Через UI GitHub:

1. Перейдите на страницу репозитория в GitHub
2. Нажмите вкладку "Pull requests"
3. Нажмите кнопку "New pull request"
4. В выпадающем списке "base:" выберите `back-dev`
5. В выпадающем списке "compare:" выберите `master`
6. Нажмите "Create pull request"
7. Заполните заголовок: `#cf Синхронизация back-dev с master. Ref CF-XX`
8. Нажмите "Create pull request"
9. Сразу нажмите "Merge pull request" и подтвердите слияние

## Разрешение конфликтов

При возникновении конфликтов:

1. Обсудите их с командой
2. Решите конфликты локально:
   ```bash
   git pull origin master
   # Разрешите конфликты в своем редакторе
   git add .
   git commit -m "#cf Разрешены конфликты с master. Ref CF-42"
   ```
3. Обновите Pull Request

Через UI GitHub:

1. Когда GitHub сообщает о конфликтах в Pull Request, нажмите кнопку "Resolve conflicts"
2. GitHub отобразит конфликтующие файлы с маркерами конфликтов
3. Отредактируйте каждый файл, удалив маркеры конфликтов (`<<<<<<<`, `=======`, `>>>>>>>`) и оставив нужный код
4. После редактирования каждого файла нажмите "Mark as resolved"
5. Когда все конфликты разрешены, нажмите "Commit merge"
6. Добавьте сообщение коммита в формате: `#cf Разрешены конфликты с master. Ref CF-42`
7. Нажмите "Commit merge" для сохранения изменений

## Визуальная схема процесса

```
master:    A---B---C----------------F-----------------I
                 \                /                 /
front-dev:        \-D---E--------+-------L--------+
                      \     \            /
front-feature:         \     G---H---J--
                        \
back-dev:                K-----------------M-------N
                                            \     /
back-feature:                                O---P
```

Где:

- A, B, C - исходные коммиты в master
- D, E - базовые коммиты в ветке front-dev
- G, H, J - коммиты в front-feature ветке
- K - базовый коммит в ветке back-dev
- O, P - коммиты в back-feature ветке
- L - коммит слияния front-feature ветки в front-dev
- M, N - коммиты в back-dev (включая слияние back-feature веток)
- F - коммит слияния front-dev в master
- I - коммит слияния back-dev в master

## Дополнительные рекомендации

- Используйте содержательные сообщения коммитов
- Делайте регулярные коммиты небольшого размера
- Всегда проверяйте, что вы находитесь в правильной ветке перед внесением изменений
- Регулярно синхронизируйте рабочие ветки с master
- Общайтесь с командой о текущих изменениях
