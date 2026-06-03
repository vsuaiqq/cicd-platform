export const helpDocs = {
  overview: {
    title: 'Справочник .cicd.yaml',
    desc:
      'Разместите файл `.cicd.yaml` в корне репозитория, чтобы описать пайплайн. Flow читает его при push-webhook, если ветка совпадает с фильтром `branches`, и выполняет джобы в порядке зависимостей.',
    tip:
      'Файл из репозитория можно переопределить в Настройки → Конфигурация пайплайна. YAML из интерфейса имеет приоритет над `.cicd.yaml` в репо.',
    minimalExample: 'Минимальный рабочий пример:',
  },
  global: {
    title: 'Глобальные поля',
    desc: 'Поддерживаемые ключи верхнего уровня в `.cicd.yaml`.',
    fields: {
      name: 'Человекочитаемое имя пайплайна. Отображается в деталях запуска и логах.',
      on: 'Конфигурация триггеров. Сейчас поддерживается только `on.push` — см. раздел триггеров.',
      env:
        'Пары ключ/значение как переменные окружения во все джобы и шаги. Переопределяется `env` джоба. Секреты и переменные проекта из настроек подмешиваются отдельно (см. Секреты и env).',
      jobs: 'Словарь джобов. Ключи — внутренние ID (буквы, цифры, дефисы).',
    },
  },
  on: {
    title: 'on — триггеры',
    desc:
      'Определяет, когда запускается пайплайн. Сейчас реализован только `push`. Другие типы (`pull_request`, `manual`, `schedule`) и path-фильтры не парсятся — если они есть в YAML, игнорируются.',
    fields: {
      push:
        'Запуск по push-webhook от подключённого git-хоста. `branches` ограничивает, с каких веток стартует run.',
      branches:
        'Список веток (`main`, `develop`, `*`). Если не указан или пуст — run на любой ветке. Только точное совпадение — glob вроде `release/*` пока не раскрываются.',
    },
    exampleNote:
      'Push в ветку, не указанную в `branches`, не запускает пайплайн.',
  },
  jobs: {
    title: 'jobs',
    desc:
      'Карта определений джобов. Ключ — уникальный ID для `needs`. Джобы без зависимостей идут параллельно; с `needs` — после успешного завершения указанных (упавший или skipped upstream блокирует зависимые).',
    tip:
      'Джоб бывает runner-джобом (есть `steps`, выполняется в контейнере или на хосте) или performance gate (есть `performance_gate`, оценивается оркестратором — без шагов и контейнера).',
  },
  jobFields: {
    title: 'Поля джоба',
    desc: 'Поля внутри определения джоба. Для обычных джобов — `steps`, для quality gate — `performance_gate`, не оба сразу.',
    fields: {
      name: 'Отображаемое имя в UI. Если не указано — используется ID джоба.',
      image:
        'Docker-образ (`golang:1.25`, `node:20`, …). Необязательно — без образа шаги выполняются на хосте раннера (без Docker и cache bind-mount). Нужен для `cache`.',
      needs:
        'ID джобов, которые должны завершиться до старта. Циклы отклоняются при разборе. Для `performance_gate` здесь обязательно должен быть `source_job`.',
      env:
        'Переменные окружения джоба. Поддерживает `${{ secrets.NAME }}` в значениях.',
      timeout:
        'Максимальное время джоба. Строки длительности Go: `5m`, `1h30m`. По умолчанию — без лимита.',
      retry:
        'Дополнительные попытки при падении. `retry: 2` или `retry: { max: 2 }` — до `max + 1` попыток (первая + ретраи).',
      approval:
        'При `required`, `true`, `yes` или `1` джоб ждёт в статусе awaiting approval, пока редактор проекта не подтвердит его на странице run.',
      cache: 'Кэш между запусками (только Docker-джобы). См. раздел cache.',
      artifacts: 'Upload путей и/или download артефактов из зависимых джобов. См. artifacts.',
      steps: 'Упорядоченный список шагов runner-джоба. Обязателен, если нет `performance_gate`.',
      performance_gate:
        'Адаптивный performance quality gate — см. раздел Performance gate. Оркестратор оценивает метрики из `source_job`; шаги раннера не выполняются.',
    },
  },
  steps: {
    title: 'steps',
    desc:
      'Упорядоченный список шагов runner-джоба. Выполняются последовательно в одном контейнере (или на хосте без `image`). Падающий шаг прерывает следующие, если не задан `continue_on_error`.',
    fields: {
      name: 'Метка шага в UI и логах.',
      run:
        'Shell-скрипт. Многострочные — блочный литерал YAML (`|`). По умолчанию `/bin/sh -e`.',
      timeout: 'Лимит времени шага (`5m`, `30s`, …).',
      retry: 'Дополнительные попытки шага. До `retry + 1` попыток всего.',
      continue_on_error:
        'При `true` джоб продолжается после ненулевого exit code шага.',
    },
  },
  cache: {
    title: 'cache',
    desc:
      'Кэширует каталоги между запусками по ключу. При hit пути восстанавливаются до шагов; при miss — сохраняются после. Требует Docker `image` — host-джобы игнорируют `cache`.',
    fields: {
      key:
        'Ключ кэша. Поддерживает `${{ checksum "relative/path" }}` — SHA-256 файла (первые 16 hex-символов), вычисляется на раннере из checkout.',
      paths:
        'Пути внутри контейнера. Относительные — от `/workspace`.',
    },
    tip:
      'Ключи общие для веток на одном раннере. Для изоляции добавляйте имя ветки или коммит в ключ как обычный текст — `${{ branch }}` не поддерживается.',
  },
  artifacts: {
    title: 'artifacts',
    desc:
      'Файлы джоба (`paths`) архивируются после успешных шагов. Другие джобы того же run могут забрать их через `download` до своих шагов. Артефакты только в рамках одного run — межrun-хранилища и кнопки скачивания в UI пока нет.',
    fields: {
      paths: 'Пути для архивации после успешного джоба. Относительно корня workspace.',
      download:
        'Список `{ job: <jobId> }`. Артефакты этих джобов распаковываются в workspace до шагов.',
    },
  },
  loadTesting: {
    title: 'Нагрузочное тестирование',
    desc:
      'Load test — обычный runner-джоб с `steps`. Специального типа джоба нет. После завершения раннер ищет `.flow/perf-metrics.json` в workspace и прикрепляет метрики к результату джоба.',
    fields: {
      metricsFile:
        'Путь: `.flow/perf-metrics.json` (фиксированный). Раннер читает файл автоматически; в YAML путь не настраивается.',
      format:
        'JSON: `version` (опционально, по умолчанию 1), опциональный `tool`, и карта `metrics` с числовыми значениями.',
      exampleMetrics:
        'Типичные ключи: `http_req_duration_p95`, `http_req_duration_avg`, `http_reqs`, `http_req_failed_rate`. Имена произвольные — gate ссылается на них по `name`.',
    },
    tip:
      'Load test работает и без performance gate. Добавьте `performance_gate` downstream, когда нужны автоматические pass/fail решения по метрикам.',
    exampleNote: 'Рабочий пример — `test-repo/scripts/load-test.sh` в репозитории платформы.',
  },
  performanceGate: {
    title: 'Performance gate',
    desc:
      'Performance gate — тип джоба платформы: сравнивает метрики load test из source-джоба с постоянными и/или адаптивными порогами. Выполняется в оркестраторе (через analytics-service) — без контейнера и шагов. При провале downstream-джобы пропускаются.',
    fields: {
      source_job:
        'ID runner-джоба, который записал `.flow/perf-metrics.json`. Должен быть в `needs` этого джоба. Нельзя указать другой performance gate.',
      metrics:
        'Список правил. У каждого: `name`, `direction` (`lower_is_better` или `higher_is_better`), опционально `max` (lower) или `min` (higher) как постоянный порог.',
      baseline_window_days:
        'Сколько дней истории учитывать для адаптивных порогов. По умолчанию: `30`.',
      baseline_min_samples:
        'Минимум исторических замеров на метрику до включения adaptive. По умолчанию: `3`. Меньше — cold start: adaptive пропускается, constant всё равно проверяется.',
      baseline_branch:
        'Опциональный фильтр ветки для baseline. По умолчанию — ветка текущего run.',
      adaptive_enabled:
        'При `true` (по умолчанию) пороги считаются по mean и stddev baseline. `false` — только constant.',
      adaptive_sigma_factor:
        'Множитель σ в adaptive-формуле. По умолчанию: `2.0`.',
      adaptive_max_regression_pct:
        'Максимальная регрессия относительно mean baseline (%). По умолчанию: `15`.',
    },
    adaptiveFormula:
      'lower_is_better: порог = min(μ + k·σ, μ × (1 + max_regression% / 100)). higher_is_better: порог = max(μ − k·σ, μ × (1 − max_regression% / 100)). Метрика проходит, если укладывается и в constant (если задан), и в эффективный порог.',
    coldStart:
      'При cold start (меньше `min_samples` исторических значений) adaptive пропускается, действуют только constant `max`/`min`. На странице run показывается бейдж cold start.',
    uiNote:
      'Результат gate на странице run: вердикт по метрикам, baseline, constant vs adaptive пороги, общий pass/fail.',
    defaultMetrics:
      'Если `metrics` не указан: `http_req_duration_p95` (lower), `http_req_failed_rate` (lower), `http_reqs` (higher) — без constant bounds.',
  },
  secrets: {
    title: 'Секреты и переменные окружения',
    desc: 'Приоритет переменных окружения (выше — важнее):',
    priorityList: [
      'Секреты — Настройки → Секреты. Наивысший приоритет, маскируются в логах.',
      'Job env — `env` в джобе.',
      'Pipeline env — верхний `env` в `.cicd.yaml`.',
      'Project env — Настройки → Переменные окружения.',
    ],
    interpolationTitle: 'Поддерживаемые подстановки',
    interpolationList: [
      '`${{ secrets.NAME }}` — только в значениях `env` пайплайна или джоба. Резолвится оркестратором до dispatch.',
      '`${{ checksum "path" }}` — только в `cache.key`. Резолвится на раннере из checkout.',
      '`$NAME` — в shell `run` для любой переменной или секрета, уже попавшего в env джоба.',
    ],
    warn:
      'Секреты только для записи после сохранения. `${{ branch }}`, `${{ sha }}`, `${{ matrix.* }}` и `if:` на джобах/шагах пока не поддерживаются.',
  },
  limitations: {
    title: 'Ограничения и roadmap',
    desc: 'Текущее поведение платформы и ещё не реализованные возможности:',
    items: [
      'Триггеры: только `on.push` с опциональным `branches`. Нет pull_request, manual-блока, schedule и path-фильтров.',
      'Условия: нет `if:` на джобах и шагах.',
      'Matrix: нет `strategy.matrix`.',
      'Environments: нет блока `environments` и поля `environment:` у джоба.',
      'allow_failure: не поддерживается — упавший джоб валит run (кроме `continue_on_error` на отдельных шагах).',
      'Step env: `env` на уровне шага не парсится из YAML.',
      'Кэш только при указанном Docker `image`. Host-джобы кэш не используют.',
      'Артефакты только в рамках одного run (`artifacts.download`). Кнопки скачивания в UI пока нет.',
      'Логи: UI показывает вывод шагов после завершения джоба; live streaming только на backend.',
      'Performance gate cold start: adaptive нуждается в истории run на той же ветке; до этого работают только constant `max`/`min`.',
      'Файл метрик: фиксированный путь `.flow/perf-metrics.json`; свой путь не настраивается.',
    ],
  },
  example: {
    title: 'Полный пример',
    desc:
      'Go-проект: lint, test, build, deploy staging, load test, adaptive performance gate, deploy production. Соответствует `test-repo/.cicd.yaml` в репозитории платформы.',
  },
} as const
