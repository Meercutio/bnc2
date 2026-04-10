    const HTTP_BASE = location.origin;
    const WS_BASE = (location.protocol === "https:" ? "wss://" : "ws://") + location.host;

    const $ = (id) => document.getElementById(id);

    // ---------------- i18n (RU/EN)
    const I18N = {
        ru: {
            lblUser: "Пользователь",
            lblTheme: "Тема",
            btnRules: "Правила",
            btnRating: "Рейтинг",
            btnAccount: "Аккаунт",
            btnClose: "Закрыть",
            btnRefresh: "Обновить",
            btnLogin: "Войти",
            btnRegister: "Регистрация",
            btnLogout: "Выйти",
            lblEmail: "Email",
            lblPassword: "Пароль",
            lblDisplayName: "Имя (для регистрации)",
            lblSession: "Сессия",
            lblAuth: "Авторизация",
            hAuth: "Вход / Регистрация",
            authAfterHtml: "После успешного входа панель переедет в меню <b>Аккаунт</b>, и откроется интерфейс игры.",
            authSecureHint: "Сессия хранится в защищённой HttpOnly cookie и не лежит в localStorage.",
            hMatch: "Матч",
            hProfile: "Профиль",
            hGame: "Игра",
            sumLogs: "Логи",
            btnCreateJoin: "Создать и подключиться",
            tipCreateJoin: "Создать матч и сразу подключиться",
            btnFindMatch: "Найти игру",
            btnCancel: "Отмена",
            btnReconnect: "Переподключиться",
            tipReconnect: "Переподключиться к последнему матчу",
            btnLeave: "Выйти из матча",
            lblMatchId: "ID матча",
            phMatchId: "вставьте id матча",
            btnJoin: "Подключиться",
            btnCopyLink: "Скопировать invite link",
            lblInviteLink: "Invite link",
            phInviteLink: "появится после id матча",
            inviteHintHtml: "Поделись ссылкой — второй игрок просто откроет её и нажмёт <b>Подключиться</b>.",
            lblYou: "Вы",
            lblPhase: "Фаза",
            lblMode: "Режим",
            lblRound: "Раунд",
            lblDeadline: "Дедлайн",
            lblSeries: "Серия",
            friendlyNote: "Дружеский матч: рейтинг не изменяется.",
            modeNoteFriendly: "Дружеский матч: игроки подключились по invite link, результат не влияет на рейтинг.",
            modeNoteRanked: "Рейтинговый матч: игроки найдены через кнопку «Найти игру», результат влияет на рейтинг.",
            lblName: "Имя",
            lblWins: "Победы",
            lblLosses: "Поражения",
            lblDraws: "Ничьи",
            lblRating: "Рейтинг",
            lblGames: "Игр",
            lblRank: "Место",
            lblWLD: "В/П/Н",
            youName: "Вы",
            oppName: "Соперник",
            statusConn: "Подключение",
            statusSecret: "Секрет",
            statusMove: "Ход",
            revealTitle: "Секреты после игры",
            lblSetSecret: "Задать секрет (4 цифры)",
            phSecret: "например 0011",
            btnSetSecret: "Задать секрет",
            lblSubmitGuess: "Отправить попытку (4 цифры)",
            phGuess: "например 0101",
            btnSubmitGuess: "Отправить",
            btnRematch: "Реванш",
            tipRematch: "Запросить реванш (нужны оба игрока)",
            lblHistory: "История",
            thRound: "Раунд",
            thWinner: "Победитель",
            thGuess: "Ход",
            thBulls: "Быки",
            thCows: "Коровы",
            hLeaderboard: "Рейтинг игроков",
            thPlayer: "Игрок",
            thRating: "Рейтинг",
            thGames: "Игр",
            rulesSub: "Как играть в Bulls & Cows",
            hRules: "Правила",
            rulesHtml: `
                <ol style="margin:0; padding-left: 18px;">
                    <li><b>Каждый игрок загадывает секрет</b> — 4 цифры (0–9). Допускаются ведущие нули и повторы.</li>
                    <li>Игроки делают попытки угадать секрет соперника. После каждого раунда вы получаете подсказку:
                        <ul style="margin:6px 0 0 0; padding-left: 18px;">
                            <li><b>Bulls</b> — цифра угадана и стоит на правильной позиции.</li>
                            <li><b>Cows</b> — цифра есть в секрете, но позиция другая.</li>
                        </ul>
                    </li>
                    <li><b>Побеждает</b> тот, кто первым угадает секрет соперника (4 Bulls).</li>
                    <li>Если включён таймер раунда и время вышло — попытка считается пропущенной.</li>
                    <li><b>Рейтинговый матч</b> начинается, когда оба игрока находят друг друга через кнопку <b>Найти игру</b>, и влияет на рейтинг.</li>
                    <li><b>Дружеский матч</b> начинается, когда игроки подключаются по invite link, и <b>не влияет</b> на рейтинг.</li>
                </ol>
            `,
            stageConnectTitle: "Подключитесь к матчу",
            stageConnectHint: "Создайте матч или вставьте ID матча, затем нажмите Подключиться.",
            stageNoConnTitle: "Нет соединения",
            stageNoConnHint: "Проверьте сеть. Мы попробуем переподключиться автоматически.",
            stageWaitOpponentTitle: "Ожидаем соперника",
            stageWaitOpponentHint: "Скопируйте invite link и отправьте второму игроку.",
            waitOpponentTitle: "Ожидаем соперника",
            waitOpponentHint: "Отправьте invite link — второй игрок откроет ссылку и нажмёт Подключиться.",
            stageSetSecretTitle: "Задайте секрет",
            stageSetSecretHint: "Введите 4 цифры (можно с нулями и повторами) и нажмите Задать секрет.",
            stageSecretOkTitle: "Секрет принят",
            stageSecretOkHint: "Ждём, пока соперник задаст секрет…",
            waitStartTitle: "Стартуем",
            waitStartHint: "Подготовка раунда…",
            stageAlmostTitle: "Почти готово",
            stageAlmostHint: "Оба секрета заданы — стартуем раунд…",
            stageMakeMoveTitle: "Сделайте ход",
            stageMakeMoveHintNoDeadline: "Введите 4 цифры и отправьте ход.",
            stageMakeMoveHintDeadline: (ts) => `Введите 4 цифры и отправьте до ${ts}.`,
            stageMoveOkTitle: "Ход принят",
            stageMoveOkHint: "Ждём ход соперника…",
            stageCountingTitle: "Раунд завершается",
            stageCountingHint: "Подсчёт результатов…",
            stageFinishedTitle: "Игра завершена",
            stageFinishedHint: "Можно запросить реванш.",
            stageFinishedHintWinner: (w) => `Победитель: ${w}. Можно запросить реванш.`,
            waitTitleDefault: "Ожидание…",
            waitLostTitle: "Потеряно соединение",
            waitLostHint: "Пробуем переподключиться…",
            netOffline: "offline",
            netOnline: "online",
            netConnecting: "подключение",
            badgeReady: "готов",
            badgeNotSet: "не задан",
            badgeWaiting: "ожидание",
            badgeNotInMatch: "не в матче",
            badgeAccepted: "принят",
            badgePending: "ожидается",
            dash: "—",
            modeRanked: "рейтинговый",
            modeFriendly: "дружеский",
            phase_waiting_players: "ожидание игроков",
            phase_waiting_secrets: "ожидание секретов",
            phase_playing: "игра",
            phase_finished: "завершено",
            mm_idle: "idle",
            mm_searching: "поиск…",
            mm_matched: "найден",
            mm_error: "ошибка",
            leaderboardTop: (n) => `Топ ${n} игроков по рейтингу (Elo)`,
            loading: "Загрузка…",
            noData: "Пока нет данных",
            failedLeaderboard: "Не удалось загрузить рейтинг",
            btnLang: "EN",
            missed: "пропуск",
            sessionGuest: "гость",
            sessionActive: "активна",
            seriesDrawLabel: "ничьи",
            authPleaseLogin: "Пожалуйста, войдите.",
            authTokenRestored: "Сессия восстановлена.",
            authSessionExpired: "Сессия истекла. Войдите снова.",
            authLoggedOut: "Вы вышли.",
            authRegistered: "Аккаунт создан. Теперь выполните вход.",
            authLoginSuccess: "Вход выполнен.",
            reconnectFailed: "Не удалось подключиться. Нажмите Переподключиться.",
            reconnectLostTitle: "Соединение потеряно",
            reconnectLostHint: (n) => `Переподключаемся… (попытка ${n})`,
            reconnectLostBanner: (reason, n) => `Connection lost${reason ? ": " + reason : ""}. Reconnecting… (attempt ${n})`,
            reconnectConnecting: (n) => `Connecting… (attempt ${n})`,
            connectMatchBanner: "Подключаемся к матчу…",
        },
        en: {
            lblUser: "User",
            lblTheme: "Theme",
            btnRules: "Rules",
            btnRating: "Rating",
            btnAccount: "Account",
            btnClose: "Close",
            btnRefresh: "Refresh",
            btnLogin: "Login",
            btnRegister: "Register",
            btnLogout: "Logout",
            lblEmail: "Email",
            lblPassword: "Password",
            lblDisplayName: "Display name (for register)",
            lblSession: "Session",
            lblAuth: "Auth",
            hAuth: "Login / Register",
            authAfterHtml: "After login, this panel moves into <b>Account</b>, and the game UI becomes available.",
            authSecureHint: "The session is kept in a secure HttpOnly cookie instead of localStorage.",
            hMatch: "Match",
            hProfile: "Profile",
            hGame: "Game",
            sumLogs: "Logs",
            btnCreateJoin: "Create & Join",
            tipCreateJoin: "Create a match and connect immediately",
            btnFindMatch: "Find match",
            btnCancel: "Cancel",
            btnReconnect: "Reconnect",
            tipReconnect: "Reconnect to the last match",
            btnLeave: "Leave match",
            lblMatchId: "Match ID",
            phMatchId: "paste match id",
            btnJoin: "Join",
            btnCopyLink: "Copy invite link",
            lblInviteLink: "Invite link",
            phInviteLink: "will appear after match id",
            inviteHintHtml: "Share the link — the 2nd player opens it and clicks <b>Join</b>.",
            lblYou: "You",
            lblPhase: "Phase",
            lblMode: "Mode",
            lblRound: "Round",
            lblDeadline: "Deadline",
            lblSeries: "Series",
            friendlyNote: "Friendly match: rating is not affected.",
            modeNoteFriendly: "Friendly match: players joined via invite link, so rating is not affected.",
            modeNoteRanked: "Ranked match: players were matched through the Find match button, so rating is affected.",
            lblName: "Name",
            lblWins: "Wins",
            lblLosses: "Losses",
            lblDraws: "Draws",
            lblRating: "Rating",
            lblGames: "Games",
            lblRank: "Rank",
            lblWLD: "W/L/D",
            youName: "You",
            oppName: "Opponent",
            statusConn: "Connection",
            statusSecret: "Secret",
            statusMove: "Move",
            revealTitle: "Secrets after game",
            lblSetSecret: "Set secret (4 digits)",
            phSecret: "e.g. 0011",
            btnSetSecret: "Set secret",
            lblSubmitGuess: "Submit guess (4 digits)",
            phGuess: "e.g. 0101",
            btnSubmitGuess: "Submit guess",
            btnRematch: "Rematch",
            tipRematch: "Request rematch (needs both players)",
            lblHistory: "History",
            thRound: "Round",
            thWinner: "Winner",
            thGuess: "Guess",
            thBulls: "Bulls",
            thCows: "Cows",
            hLeaderboard: "Leaderboard",
            thPlayer: "Player",
            thRating: "Rating",
            thGames: "Games",
            rulesSub: "How to play Bulls & Cows",
            hRules: "Rules",
            rulesHtml: `
                <ol style="margin:0; padding-left: 18px;">
                    <li><b>Each player sets a secret</b> — 4 digits (0–9). Leading zeros and repeats are allowed.</li>
                    <li>Players submit guesses. After each round you get hints:
                        <ul style="margin:6px 0 0 0; padding-left: 18px;">
                            <li><b>Bulls</b> — correct digit in the correct position.</li>
                            <li><b>Cows</b> — digit exists, but in a different position.</li>
                        </ul>
                    </li>
                    <li><b>You win</b> by guessing the opponent’s secret first (4 Bulls).</li>
                    <li>If the round timer is enabled and time runs out — your attempt is missed.</li>
                    <li><b>Ranked match</b> starts when both players find each other through <b>Find match</b>, and it affects rating.</li>
                    <li><b>Friendly match</b> starts when players join via invite link, and it <b>does not</b> affect rating.</li>
                </ol>
            `,
            stageConnectTitle: "Connect to a match",
            stageConnectHint: "Create a match or paste Match ID, then click Join.",
            stageNoConnTitle: "No connection",
            stageNoConnHint: "Check your network. We'll try to reconnect automatically.",
            stageWaitOpponentTitle: "Waiting for opponent",
            stageWaitOpponentHint: "Copy invite link and send it to the 2nd player.",
            waitOpponentTitle: "Waiting for opponent",
            waitOpponentHint: "Send invite link — the 2nd player opens it and clicks Join.",
            stageSetSecretTitle: "Set your secret",
            stageSetSecretHint: "Enter 4 digits (zeros and repeats allowed) and click Set secret.",
            stageSecretOkTitle: "Secret accepted",
            stageSecretOkHint: "Waiting for opponent to set a secret…",
            waitStartTitle: "Starting",
            waitStartHint: "Preparing the round…",
            stageAlmostTitle: "Almost ready",
            stageAlmostHint: "Both secrets are set — starting…",
            stageMakeMoveTitle: "Make your move",
            stageMakeMoveHintNoDeadline: "Enter 4 digits and submit.",
            stageMakeMoveHintDeadline: (ts) => `Enter 4 digits and submit before ${ts}.`,
            stageMoveOkTitle: "Move accepted",
            stageMoveOkHint: "Waiting for opponent’s move…",
            stageCountingTitle: "Finishing round",
            stageCountingHint: "Calculating results…",
            stageFinishedTitle: "Game finished",
            stageFinishedHint: "You can request a rematch.",
            stageFinishedHintWinner: (w) => `Winner: ${w}. You can request a rematch.`,
            waitTitleDefault: "Waiting…",
            waitLostTitle: "Connection lost",
            waitLostHint: "Reconnecting…",
            netOffline: "offline",
            netOnline: "online",
            netConnecting: "connecting",
            badgeReady: "ready",
            badgeNotSet: "not set",
            badgeWaiting: "waiting",
            badgeNotInMatch: "not in match",
            badgeAccepted: "accepted",
            badgePending: "pending",
            dash: "—",
            modeRanked: "ranked",
            modeFriendly: "friendly",
            phase_waiting_players: "waiting players",
            phase_waiting_secrets: "waiting secrets",
            phase_playing: "playing",
            phase_finished: "finished",
            mm_idle: "idle",
            mm_searching: "searching…",
            mm_matched: "matched",
            mm_error: "error",
            leaderboardTop: (n) => `Top ${n} players by rating (Elo)`,
            loading: "Loading…",
            noData: "No data yet",
            failedLeaderboard: "Failed to load leaderboard",
            btnLang: "RU",
            missed: "missed",
            sessionGuest: "guest",
            sessionActive: "active",
            seriesDrawLabel: "draw",
            authPleaseLogin: "Please login.",
            authTokenRestored: "Session restored.",
            authSessionExpired: "Session expired. Please login again.",
            authLoggedOut: "Logged out.",
            authRegistered: "Account created. Now log in.",
            authLoginSuccess: "Logged in.",
            reconnectFailed: "Failed to connect. Click Reconnect.",
            reconnectLostTitle: "Connection lost",
            reconnectLostHint: (n) => `Reconnecting… (attempt ${n})`,
            reconnectLostBanner: (reason, n) => `Connection lost${reason ? ": " + reason : ""}. Reconnecting… (attempt ${n})`,
            reconnectConnecting: (n) => `Connecting… (attempt ${n})`,
            connectMatchBanner: "Connecting to match…",
        }
    };

    let currentLang = localStorage.getItem("lang") || "ru";

    function t(key, ...args) {
        const dict = I18N[currentLang] || I18N.ru;
        const v = dict[key] ?? (I18N.ru[key] ?? key);
        return typeof v === "function" ? v(...args) : v;
    }

    function applyI18n() {
        document.documentElement.lang = currentLang;
        // textContent
        document.querySelectorAll("[data-i18n]").forEach((el) => {
            const k = el.getAttribute("data-i18n");
            el.textContent = t(k);
        });
        // placeholders
        document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
            const k = el.getAttribute("data-i18n-placeholder");
            el.placeholder = t(k);
        });
        // titles
        document.querySelectorAll("[data-i18n-title]").forEach((el) => {
            const k = el.getAttribute("data-i18n-title");
            el.title = t(k);
        });
        // html content
        document.querySelectorAll("[data-i18n-html]").forEach((el) => {
            const k = el.getAttribute("data-i18n-html");
            el.innerHTML = t(k);
        });

        // Theme option labels (both selects)
        document.querySelectorAll('#themeSelect option[value="system"], #themeSelectTop option[value="system"]').forEach(o => o.textContent = currentLang === "ru" ? "Системная" : "System");
        document.querySelectorAll('#themeSelect option[value="light"], #themeSelectTop option[value="light"]').forEach(o => o.textContent = currentLang === "ru" ? "Светлая" : "Light");
        document.querySelectorAll('#themeSelect option[value="dark"], #themeSelectTop option[value="dark"]').forEach(o => o.textContent = currentLang === "ru" ? "Тёмная" : "Dark");

        // Language button label
        const btn = $("btnLang");
        if (btn) btn.textContent = t("btnLang");
        setSessionStatus(sessionActive);
        setTableHeaders(lastNames.p1 || "p1", lastNames.p2 || "p2");

        // If we already have state — re-render dynamic strings in the current language
        // Matchmaking status label
        try { if (typeof mmStatusRaw !== "undefined") setMMStatus(mmStatusRaw); } catch {}

        if (lastState) updateUXFromState(lastState);
        else {
            renderSeries({ p1Wins: 0, p2Wins: 0, draws: 0 }, lastNames);
            if ($("historyBody")) {
                $("historyBody").innerHTML = `<tr><td colspan="8" class="muted">${t("noData")}</td></tr>`;
            }
            setStageCopy(t("stageConnectTitle"), t("stageConnectHint"));
            resetActionPanels();
        }
    }

    function toggleLang() {
        currentLang = currentLang === "ru" ? "en" : "ru";
        localStorage.setItem("lang", currentLang);
        applyI18n();
    }

    // ---------------- Theme
    function applyTheme(v) {
        const html = document.documentElement;
        html.setAttribute("data-theme", v);
        localStorage.setItem("theme", v);
        if ($("themeSelect")) $("themeSelect").value = v;
        if ($("themeSelectTop")) $("themeSelectTop").value = v;
    }
    function initTheme() {
        const v = localStorage.getItem("theme") || "system";
        applyTheme(v);
        if ($("themeSelect")) $("themeSelect").onchange = (e) => applyTheme(e.target.value);
        if ($("themeSelectTop")) $("themeSelectTop").onchange = (e) => applyTheme(e.target.value);
    }

    // ---------------- Session
    let sessionActive = false;

    function setSessionStatus(active) {
        sessionActive = !!active;
        const el = $("sessionStatus");
        if (el) el.textContent = active ? t("sessionActive") : t("sessionGuest");
    }

    function setAuthMsg(text, isErr=false) {
        const el = $("authMsg");
        if (!el) return;
        el.textContent = text || "";
        el.className = isErr ? "err" : "ok";
    }

    function log(line) {
        const el = $("log");
        if (!el) return;
        el.textContent += line + "\n";
        el.scrollTop = el.scrollHeight;
    }

    function formatError(err) {
        if (err?.body?.message) return err.body.message;
        if (typeof err?.body === "string" && err.body.trim()) return err.body;
        if (err?.message) return err.message;
        try { return JSON.stringify(err); } catch {}
        return String(err || "unknown error");
    }

    function handleUnauthorizedSession() {
        disconnectWS();
        setSessionStatus(false);
        setAuthedUI(false);
        setAuthMsg(t("authSessionExpired"), true);
    }

    // ---------------- API
    async function api(path, opts={}) {
        const headers = { ...(opts.headers || {}) };
        if (opts.body !== undefined && !headers["Content-Type"]) {
            headers["Content-Type"] = "application/json";
        }

        const res = await fetch(HTTP_BASE + path, {
            ...opts,
            headers,
            credentials: "same-origin",
        });
        const text = await res.text();
        let json = null;
        try { json = text ? JSON.parse(text) : null; } catch {}
        if (!res.ok) {
            if (res.status === 401 && path !== "/api/auth/login") {
                handleUnauthorizedSession();
            }
            throw { status: res.status, body: json || text };
        }
        return json;
    }

    function syncAuthPanel(isAuthed) {
        $("authFormFields")?.classList.toggle("hidden", !!isAuthed);
        $("authGuestActions")?.classList.toggle("hidden", !!isAuthed);
        $("authSessionActions")?.classList.toggle("hidden", !isAuthed);
        $("authSecureHint")?.classList.toggle("hidden", !isAuthed);
    }

    // ---------------- UI mode
    function setAuthedUI(isAuthed) {
        const overlay = $("overlay");
        const app = $("app");
        const authPanel = $("authPanel");
        const overlaySlot = $("overlaySlot");
        const headerAuthSlot = $("headerAuthSlot");
        const accountDetails = $("accountDetails");

        if (!overlay || !app || !authPanel || !overlaySlot || !headerAuthSlot) return;
        syncAuthPanel(isAuthed);
        setSessionStatus(isAuthed);

        if (isAuthed) {
            overlay.classList.add("hidden");
            app.classList.remove("hidden");
            headerAuthSlot.appendChild(authPanel);
            authPanel.classList.add("compact");
            if (accountDetails) accountDetails.open = false;
        } else {
            app.classList.add("hidden");
            overlay.classList.remove("hidden");
            overlaySlot.appendChild(authPanel);
            authPanel.classList.remove("compact");
            if ($("userBadge")) $("userBadge").textContent = "-";
        }
    }

    async function refreshMe() {
        const me = await api("/api/me", { method: "GET" });
        const name = (me?.displayName || "").trim() || "-";
        const email = (me?.email || "").trim() || "-";
        const st = me?.stats || {};

        setSessionStatus(true);

        window.__meId = me?.id || "";
        window.__meName = name;

        $("userBadge").textContent = `${name} (${email})`;
        $("profName").textContent = name;
        $("profEmail").textContent = email;
        $("profWins").textContent = String(st.wins ?? 0);
        $("profLosses").textContent = String(st.losses ?? 0);
        $("profDraws").textContent = String(st.draws ?? 0);

        $("profRating").textContent = String(st.rating ?? 1000);
        $("profGames").textContent = String(st.games ?? 0);

        // rank is computed via window-function, keep it separate
        try {
            const r = await api("/api/rating/me", { method: "GET" });
            $("profRank").textContent = (r?.rank ? String(r.rank) : "-");
        } catch {
            $("profRank").textContent = "-";
        }
        return me;
    }

    // ---------------- Rating / leaderboard
    function openRating() {
        $("ratingOverlay").classList.remove("hidden");
        loadRating();
    }

    function closeRating() {
        $("ratingOverlay").classList.add("hidden");
    }

    async function loadRating() {
        const LIMIT = 20;
        const body = $("ratingBody");
        body.innerHTML = `<tr><td colspan="7" class="muted">${t("loading")}</td></tr>`;
        $("ratingSub").textContent = t("leaderboardTop", LIMIT);

        // my bar
        $("myName").textContent = window.__meName || "-";
        $("myRating").textContent = "-";
        $("myRank").textContent = "-";
        $("myWLD").textContent = "-";

        try {
            const lb = await api(`/api/rating/leaderboard?limit=${LIMIT}`, { method: "GET" });
            const items = (lb?.items || []).slice(0, LIMIT);

            // Best-effort my rating
            try {
                const me = await api("/api/rating/me", { method: "GET" });
                if (me) {
                    $("myName").textContent = me.displayName || (window.__meName || "-");
                    $("myRating").textContent = String(me.rating ?? "-");
                    $("myRank").textContent = me.rank ? String(me.rank) : "-";
                    $("myWLD").textContent = `${me.wins ?? 0}/${me.losses ?? 0}/${me.draws ?? 0}`;
                }
            } catch {}

            if (!items.length) {
                body.innerHTML = `<tr><td colspan="7" class="muted">${t("noData")}</td></tr>`;
                return;
            }

            body.innerHTML = items.map((e) => {
                const isMe = window.__meId && e.userId === window.__meId;
                return `
                    <tr class="${isMe ? "hlRow" : ""}">
                        <td>${e.rank}</td>
                        <td>${escapeHtml(e.displayName || "-")}</td>
                        <td>${e.rating}</td>
                        <td>${e.games}</td>
                        <td>${e.wins}</td>
                        <td>${e.losses}</td>
                        <td>${e.draws}</td>
                    </tr>
                `;
            }).join("");
        } catch (e) {
            body.innerHTML = `<tr><td colspan="7" class="err">${t("failedLeaderboard")}</td></tr>`;
        }
    }


    // ---------------- Rules modal
    function openRules() {
        $("rulesOverlay").classList.remove("hidden");
    }
    function closeRules() {
        $("rulesOverlay").classList.add("hidden");
    }
    // ---------------- MatchId + invite link
    function getQueryMatchId() {
        const u = new URL(location.href);
        const m = (u.searchParams.get("match") || "").trim();
        return m;
    }

    function setInviteLink(matchId) {
        const u = new URL(location.href);
        if (matchId) u.searchParams.set("match", matchId);
        else u.searchParams.delete("match");
        $("inviteLink").value = matchId ? u.toString() : "";
    }

    function setMatchId(matchId, {persist=true, updateUrl=true} = {}) {
        $("matchId").value = matchId || "";
        setInviteLink(matchId || "");
        if (persist) {
            if (matchId) localStorage.setItem("lastMatchId", matchId);
            else localStorage.removeItem("lastMatchId");
        }
        if (updateUrl) {
            const u = new URL(location.href);
            if (matchId) u.searchParams.set("match", matchId);
            else u.searchParams.delete("match");
            history.replaceState(null, "", u.toString());
        }
    }

    let mmStatusRaw = "idle";
    function setMMStatus(raw) {
        mmStatusRaw = raw || "idle";
        document.getElementById("mmStatus").textContent = mmLabel(mmStatusRaw);
    }

    let mmAbort = null;

    function setMMUI(active) {
        document.getElementById("btnFindMatch").disabled = active;
        document.getElementById("btnCancelFind").classList.toggle("hidden", !active);
    }

    async function matchmakingJoin(signal) {
        const res = await fetch(HTTP_BASE + "/api/matchmaking/join", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            signal,
            credentials: "same-origin",
        });
        const text = await res.text();
        let json = null;
        try { json = text ? JSON.parse(text) : null; } catch {}
        return { status: res.status, json };
    }

    async function findMatch() {
        if (mmAbort) return;

        disconnectWS(); // если уже были подключены к матчу

        setMMStatus("searching");
        setMMUI(true);

        mmAbort = new AbortController();
        try {
            const out = await matchmakingJoin(mmAbort.signal);
            if (out.status === 200 && out.json?.matchId) {
                const m = out.json.matchId;
                setMMStatus("matched");
                setMatchId(m, {persist:true, updateUrl:true});
                connectWS(m);
                log("[mm] matched: " + m);
            } else if (out.status === 401) {
                handleUnauthorizedSession();
                setMMStatus("error");
            } else if (out.status === 409) {
                setMMStatus("searching"); // уже ждём
            } else {
                setMMStatus("error");
                log("[mm] unexpected: " + JSON.stringify(out));
            }
        } catch (e) {
            if (e?.name === "AbortError") {
                setMMStatus("idle");
                log("[mm] canceled");
            } else {
                setMMStatus("error");
                log("[mm] error: " + JSON.stringify(e));
            }
        } finally {
            mmAbort = null;
            setMMUI(false);
            if (mmStatusRaw === "searching") setMMStatus("idle");
        }
    }

    function cancelFindMatch() {
        if (!mmAbort) return;
        mmAbort.abort();
    }

    document.getElementById("btnFindMatch").onclick = () => findMatch();
    document.getElementById("btnCancelFind").onclick = () => cancelFindMatch();


    // ---------------- WS
    let ws = null;
    let currentMatchId = "";
    let wantReconnect = false;
    let reconnectAttempt = 0;
    let reconnectTimer = null;
    let lastState = null;
    let wsEverOpened = false;

    let lastNames = { p1: "p1", p2: "p2" };

    function normalizeNames(names) {
        const p1 = (names?.p1 || "").trim() || "p1";
        const p2 = (names?.p2 || "").trim() || "p2";
        return { p1, p2 };
    }

    function setTableHeaders(p1, p2) {
        $("hP1Guess").textContent = `${p1} ${t("thGuess")}`;
        $("hP1Bulls").textContent = `${p1} ${t("thBulls")}`;
        $("hP1Cows").textContent = `${p1} ${t("thCows")}`;
        $("hP2Guess").textContent = `${p2} ${t("thGuess")}`;
        $("hP2Bulls").textContent = `${p2} ${t("thBulls")}`;
        $("hP2Cows").textContent = `${p2} ${t("thCows")}`;
    }

    function winnerLabel(winner, names) {
        if (!winner) return "";
        if (winner === "draw") return "draw";
        if (winner === "p1") return names.p1;
        if (winner === "p2") return names.p2;
        return winner;
    }

    function phaseLabel(phase) {
        if (!phase) return "-";
        return t("phase_" + phase);
    }

    function winnerLabelLocalized(winner, names) {
        if (!winner) return "";
        if (winner === "draw") return currentLang === "ru" ? "ничья" : "draw";
        return winnerLabel(winner, names);
    }

    function mmLabel(raw) {
        if (!raw) return "-";
        return t("mm_" + raw);
    }

    function setWSStatus(s) { $("wsStatus").textContent = s; }

    function send(type, payload) {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            log("[ws] not connected");
            return;
        }
        ws.send(JSON.stringify({ type, payload }));
    }

    function escapeHtml(s) {
        return String(s)
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll('"', "&quot;")
            .replaceAll("'", "&#039;");
    }

    function renderHistory(state) {
        const body = $("historyBody");
        const hist = state?.history || [];
        const winner = state?.winner || "";

        const names = normalizeNames(state?.playerNames);
        lastNames = names;

        const lastIdx = hist.length - 1;

        if (!hist.length) {
            body.innerHTML = `<tr><td colspan="8" class="muted">${t("noData")}</td></tr>`;
            return;
        }

        body.innerHTML = hist.map((item, idx) => {
            const r = item.round ?? "-";
            const p1 = item.p1 || {};
            const p2 = item.p2 || {};
            const rowClass = idx === lastIdx ? "histLatest" : "";

            let p1Guess = p1.guess ?? (p1.missed ? `(${t("missed")})` : "");
            let p2Guess = p2.guess ?? (p2.missed ? `(${t("missed")})` : "");

            const p1B = p1.bulls ?? 0;
            const p1C = p1.cows ?? 0;
            const p2B = p2.bulls ?? 0;
            const p2C = p2.cows ?? 0;

            const winnerCell = (winner && idx === lastIdx) ? escapeHtml(winnerLabelLocalized(winner, names)) : "";

            return `
        <tr class="${rowClass}">
          <td>${r}</td>
          <td>${escapeHtml(p1Guess || "")}</td>
          <td>${p1B}</td>
          <td>${p1C}</td>
          <td>${escapeHtml(p2Guess || "")}</td>
          <td>${p2B}</td>
          <td>${p2C}</td>
          <td>${winnerCell}</td>
        </tr>
      `;
        }).join("");
    }

    
    // ---------------- Network UI helpers
    function setNetBanner(show, text, dotKind = "muted", showBtn = true) {
        const b = $("netBanner");
        if (!b) return;
        if (!show) {
            b.classList.add("hidden");
            return;
        }
        b.classList.remove("hidden");
        $("netText").textContent = text || "";
        const dot = $("netDot");
        dot.className = "netDot" + (dotKind === "ok" ? " ok" : dotKind === "err" ? " err" : "");
        $("btnReconnectNow").classList.toggle("hidden", !showBtn);
    }

    function setNetBadge(kind, text) {
        const el = $("netBadge");
        el.className = "badge" + (kind ? " " + kind : "");
        el.textContent = text;
    }

    function setWaitOverlay(show, title, hint, opts = {}) {
        const options = typeof opts === "boolean" ? { showReconnect: opts } : opts;
        const o = $("waitOverlay");
        if (!o) return;
        o.classList.toggle("hidden", !show);
        o.setAttribute("aria-hidden", (!show).toString());
        o.dataset.tone = options.tone || "";
        $("waitTitle").textContent = title || t("waitTitleDefault");
        $("waitHint").textContent = hint || "";
        $("btnOverlayReconnect").classList.toggle("hidden", !options.showReconnect);
        $("btnOverlayCopyLink").classList.toggle("hidden", !options.showCopy || !$("inviteLink").value.trim());
    }

    function setBadge(id, kind, text) {
        const el = $(id);
        el.className = "badge" + (kind ? " " + kind : "");
        el.textContent = text;
    }

    function otherSlot(you) { return you === "p1" ? "p2" : "p1"; }

    function renderSeries(series, names) {
        const score = series || {};
        $("series").textContent = `${names.p1} ${score.p1Wins || 0} : ${score.p2Wins || 0} ${names.p2} (${t("seriesDrawLabel")} ${score.draws || 0})`;
    }

    function setStageCopy(title, hint) {
        if ($("stageTitle")) $("stageTitle").textContent = title || t("dash");
        if ($("stageHint")) $("stageHint").textContent = hint || t("dash");
        if ($("stageLead")) $("stageLead").textContent = hint || t("dash");
    }

    function setActionPanel(kind, state, statusText, hintText) {
        const card = $(`${kind}Card`);
        const status = $(`${kind}State`);
        const hint = $(`${kind}Hint`);
        if (card) card.dataset.state = state || "waiting";
        if (status) status.textContent = statusText || t("dash");
        if (hint) hint.textContent = hintText || t("dash");
    }

    function resetActionPanels() {
        setActionPanel("secret", "waiting", t("badgeWaiting"), t("stageConnectHint"));
        setActionPanel("guess", "waiting", t("badgeWaiting"), t("stageConnectHint"));
        setActionPanel("rematch", "waiting", t("badgeWaiting"), t("tipRematch"));
    }

    function updateUXFromState(s) {
        lastState = s;

        const names = normalizeNames(s.playerNames);
        lastNames = names;

        const you = s.you || "p1";
        const opp = otherSlot(you);

        // Header pills (existing)
        setTableHeaders(names.p1, names.p2);
        $("youSlot").textContent = (names[you] || you || "-");
        $("phase").textContent = phaseLabel(s.phase);
        $("round").textContent = (s.round ?? "-");
        $("deadline").textContent = s.deadlineMs ? new Date(s.deadlineMs).toLocaleTimeString() : "-";
        renderSeries(s.series, names);

        if ($("matchMode")) {
            $("matchMode").textContent = s.ranked ? t("modeRanked") : t("modeFriendly");
        }
        if ($("modeNote")) {
            $("modeNote").textContent = s.ranked ? t("modeNoteRanked") : t("modeNoteFriendly");
            $("modeNote").classList.remove("hidden");
        }

        // Status board
        $("youName").textContent = names[you] || t("youName");
        $("oppName").textContent = names[opp] || t("oppName");
        setBadge("youRole", "muted", you.toUpperCase());
        setBadge("oppRole", "muted", opp.toUpperCase());

        // Connected (server gives only total count; we treat it as opponent online/offline)
        const oppOnline = (s.playersConnected || 0) >= 2;
        setBadge("youConn", "ok", t("netOnline"));
        setBadge("oppConn", oppOnline ? "ok" : "muted", oppOnline ? t("netOnline") : t("netOffline"));

        // Secrets
        const sr = s.secretsReady || {};
        const youSecretReady = !!sr[you];
        const oppSecretReady = !!sr[opp];
        setBadge("youSecret", youSecretReady ? "ok" : "muted", youSecretReady ? t("badgeReady") : t("badgeNotSet"));
        setBadge("oppSecret", oppSecretReady ? "ok" : "muted", oppSecretReady ? t("badgeReady") : (s.phase === "waiting_players" ? t("badgeNotInMatch") : t("badgeWaiting")));

        // Guesses (current round)
        const gr = s.guessesReady || {};
        const youGuessReady = !!gr[you];
        const oppGuessReady = !!gr[opp];
        if (s.phase === "playing") {
            setBadge("youGuess", youGuessReady ? "ok" : "muted", youGuessReady ? t("badgeAccepted") : t("badgePending"));
            setBadge("oppGuess", oppGuessReady ? "ok" : "muted", oppGuessReady ? t("badgeAccepted") : t("badgePending"));
        } else {
            setBadge("youGuess", "muted", t("dash"));
            setBadge("oppGuess", "muted", t("dash"));
        }

        // Reveal secrets (only after finished)
        const rb = $("revealBox");
        if (s.phase === "finished" && s.revealedSecrets) {
            rb.classList.remove("hidden");
            $("revealYouLabel").textContent = names[you] || "Вы";
            $("revealOppLabel").textContent = names[opp] || "Соперник";
            $("revealYouSecret").textContent = (s.revealedSecrets[you] || "----");
            $("revealOppSecret").textContent = (s.revealedSecrets[opp] || "----");
        } else {
            rb.classList.add("hidden");
        }

        // Controls enable/disable
        const wsOpen = !!ws && ws.readyState === WebSocket.OPEN;
        const setSecretEnabled =
            wsOpen && (s.phase === "waiting_secrets") && !youSecretReady;
        const guessEnabled =
            wsOpen && (s.phase === "playing") && !youGuessReady;

        $("secret").disabled = !setSecretEnabled;
        $("btnSetSecret").disabled = !setSecretEnabled;

        $("guess").disabled = !guessEnabled;
        $("btnGuess").disabled = !guessEnabled;

        const rematchEnabled = wsOpen && s.phase === "finished";
        $("btnRematch").disabled = !rematchEnabled;

        const moveHint = s.deadlineMs
            ? t("stageMakeMoveHintDeadline", new Date(s.deadlineMs).toLocaleTimeString())
            : t("stageMakeMoveHintNoDeadline");

        // Stage coach + waiting overlay
        if (!wsOpen) {
            setStageCopy(t("stageNoConnTitle"), t("stageNoConnHint"));
            setNetBadge("muted", t("netOffline"));
            setWaitOverlay(true, t("waitLostTitle"), t("waitLostHint"), { showReconnect: true, tone: "warning" });
            setActionPanel("secret", "waiting", t("badgeWaiting"), t("stageNoConnHint"));
            setActionPanel("guess", "waiting", t("badgeWaiting"), t("stageNoConnHint"));
            setActionPanel("rematch", "waiting", t("badgeWaiting"), t("stageNoConnHint"));
            return;
        }

        setNetBadge("ok", t("netOnline"));

        if (s.phase === "waiting_players") {
            setStageCopy(t("stageWaitOpponentTitle"), t("stageWaitOpponentHint"));
            setWaitOverlay(true, t("waitOpponentTitle"), t("waitOpponentHint"), { showCopy: true });
            setActionPanel("secret", "waiting", t("badgeWaiting"), t("stageWaitOpponentHint"));
            setActionPanel("guess", "waiting", t("badgeWaiting"), t("stageWaitOpponentHint"));
            setActionPanel("rematch", "waiting", t("badgeWaiting"), t("tipRematch"));
        } else if (s.phase === "waiting_secrets") {
            if (!youSecretReady) {
                setStageCopy(t("stageSetSecretTitle"), t("stageSetSecretHint"));
                setWaitOverlay(false, "", "", {});
                setActionPanel("secret", "active", t("stageSetSecretTitle"), t("stageSetSecretHint"));
                setActionPanel("guess", "waiting", t("badgeWaiting"), t("stageAlmostHint"));
                setActionPanel("rematch", "waiting", t("badgeWaiting"), t("tipRematch"));
            } else if (!oppSecretReady) {
                setStageCopy(t("stageSecretOkTitle"), t("stageSecretOkHint"));
                setWaitOverlay(true, t("waitOpponentTitle"), t("stageSecretOkHint"), {});
                setActionPanel("secret", "done", t("badgeReady"), t("stageSecretOkHint"));
                setActionPanel("guess", "waiting", t("badgeWaiting"), t("stageSecretOkHint"));
                setActionPanel("rematch", "waiting", t("badgeWaiting"), t("tipRematch"));
            } else {
                setStageCopy(t("stageAlmostTitle"), t("stageAlmostHint"));
                setWaitOverlay(true, t("waitStartTitle"), t("waitStartHint"), {});
                setActionPanel("secret", "done", t("badgeReady"), t("stageAlmostHint"));
                setActionPanel("guess", "waiting", t("badgeWaiting"), t("stageAlmostHint"));
                setActionPanel("rematch", "waiting", t("badgeWaiting"), t("tipRematch"));
            }
        } else if (s.phase === "playing") {
            setActionPanel("secret", youSecretReady ? "done" : "waiting", youSecretReady ? t("badgeReady") : t("badgeWaiting"), youSecretReady ? t("stageSecretOkHint") : t("stageSetSecretHint"));
            if (!youGuessReady) {
                setStageCopy(t("stageMakeMoveTitle"), moveHint);
                setWaitOverlay(false, "", "", {});
                setActionPanel("guess", "active", t("stageMakeMoveTitle"), moveHint);
            } else if (!oppGuessReady) {
                setStageCopy(t("stageMoveOkTitle"), t("stageMoveOkHint"));
                setWaitOverlay(true, t("waitOpponentTitle"), t("stageMoveOkHint"), {});
                setActionPanel("guess", "done", t("badgeAccepted"), t("stageMoveOkHint"));
            } else {
                setStageCopy(t("stageCountingTitle"), t("stageCountingHint"));
                setWaitOverlay(true, t("stageCountingTitle"), t("stageCountingHint"), {});
                setActionPanel("guess", "done", t("badgeAccepted"), t("stageCountingHint"));
            }
            setActionPanel("rematch", "waiting", t("badgeWaiting"), t("tipRematch"));
        } else if (s.phase === "finished") {
            const w = s.winner || "";
            const finalHint = w ? t("stageFinishedHintWinner", winnerLabelLocalized(w, names)) : t("stageFinishedHint");
            setStageCopy(t("stageFinishedTitle"), finalHint);
            setWaitOverlay(false, "", "", {});
            setActionPanel("secret", "done", t("badgeReady"), finalHint);
            setActionPanel("guess", "done", t("stageFinishedTitle"), finalHint);
            setActionPanel("rematch", "active", t("stageFinishedTitle"), finalHint);
        } else {
            setStageCopy(t("hMatch"), t("dash"));
            setWaitOverlay(false, "", "", {});
            resetActionPanels();
        }
    }

    // ---------------- Auto-reconnect
    function clearReconnect() {
        if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            reconnectTimer = null;
        }
    }

    function stopAutoReconnect() {
        wantReconnect = false;
        reconnectAttempt = 0;
        clearReconnect();
        setNetBanner(false);
        setWaitOverlay(false, "", "", false);
        setNetBadge("muted", t("netOffline"));
        wsEverOpened = false;
    }

    function scheduleReconnect(reason = "") {
        if (!wantReconnect || !currentMatchId) return;
        if (reconnectTimer) return;

        const base = 500;               // ms
        const max = 8000;               // ms
        const capAttempts = 12;
        if (reconnectAttempt >= capAttempts) {
            setWSStatus("closed");
            setNetBanner(true, t("reconnectFailed"), "err", true);
            // Не показываем "потеряно соединение" при первом заходе (когда ещё ни разу не было успешного connect).
            if (wsEverOpened) {
                setWaitOverlay(true, t("reconnectLostTitle"), t("reconnectFailed"), true);
            } else {
                setWaitOverlay(false, "", "", false);
            }
            return;
        }

        const backoff = Math.min(max, base * Math.pow(2, reconnectAttempt));
        const jitter = Math.floor(Math.random() * 250);
        const wait = backoff + jitter;

        setWSStatus("reconnecting");
        if (wsEverOpened) {
            setNetBanner(true, t("reconnectLostBanner", reason, reconnectAttempt + 1), "err", true);
            setWaitOverlay(true, t("reconnectLostTitle"), t("reconnectLostHint", reconnectAttempt + 1), true);
        } else {
            // Первый заход/авто-join: не пугаем "потеряно соединение", просто пытаемся подключиться.
            setNetBanner(true, t("reconnectConnecting", reconnectAttempt + 1), "muted", true);
            setWaitOverlay(false, "", "", false);
        }

        reconnectTimer = setTimeout(() => {
            reconnectTimer = null;
            reconnectAttempt++;
            openWS(currentMatchId, true);
        }, wait);
    }

    function disconnectWS(intentional = true) {
        clearReconnect();
        if (intentional) stopAutoReconnect();

        if (ws) {
            try { ws.close(); } catch {}
        }
        ws = null;
        setWSStatus("closed");
        if (intentional) {
            setNetBanner(false);
            setWaitOverlay(false, "", "", false);
        }
    }

    function openWS(matchId, isRetry = false) {
        if (!matchId) return log("[ui] missing matchId");
        if (!sessionActive) return log("[ui] " + t("authPleaseLogin"));

        clearReconnect();

        // Close previous socket (events from old sockets are ignored via guards below)
        const prev = ws;
        if (prev) {
            try { prev.close(); } catch {}
        }

        const url = `${WS_BASE}/ws/${encodeURIComponent(matchId)}`;
        const sock = new WebSocket(url);
        ws = sock;

        setWSStatus(isRetry ? "reconnecting" : "connecting");
        setNetBadge("muted", t("netConnecting"));
        if (!isRetry) setNetBanner(true, t("connectMatchBanner"), "muted", true);

        sock.onopen = () => {
            if (ws !== sock) return;
            setWSStatus("open");
            setNetBanner(false);
            setWaitOverlay(false, "", "", false);
            setNetBadge("ok", t("netOnline"));
            reconnectAttempt = 0;
            wsEverOpened = true;
            log("[ws] connected");

            // refresh UX with last known state (if any)
            if (lastState) updateUXFromState(lastState);
        };

        sock.onclose = (ev) => {
            if (ws !== sock) return;
            ws = null;
            setWSStatus("closed");
            setNetBadge("muted", t("netOffline"));
            const reason = (ev && ev.reason) ? ev.reason : "";
            log("[ws] closed" + (reason ? " (" + reason + ")" : ""));
            scheduleReconnect(reason);
        };

        sock.onerror = () => {
            if (ws !== sock) return;
            setWSStatus("error");
            setNetBadge("muted", t("netOffline"));
            log("[ws] error");
            // close will usually follow, but schedule just in case
            scheduleReconnect("ws error");
        };

        sock.onmessage = (ev) => {
            if (ws !== sock) return;
            let msg = null;
            try { msg = JSON.parse(ev.data); } catch { return; }

            if (msg.type === "state") {
                const s = msg.payload || {};
                updateUXFromState(s);
                renderHistory(s);
            } else if (msg.type === "series_score") {
                const ser = msg.payload?.series || {};
                renderSeries(ser, lastNames);
                log("[series_score] " + JSON.stringify(ser));
            } else if (msg.type === "game_finished") {
                const w = msg.payload?.winner || "?";
                log("[game_finished] winner=" + winnerLabel(w, lastNames));

                // рейтинг обновляется на сервере асинхронно — подтянем спустя небольшой лаг
                setTimeout(() => {
                    refreshMe().catch(() => {});
                    if (!$("ratingOverlay").classList.contains("hidden")) loadRating();
                }, 600);
            } else if (msg.type === "error") {
                log("[error] " + JSON.stringify(msg.payload));
            } else {
                log("[msg] " + JSON.stringify(msg));
            }
        };
    }

function connectWS(matchId) {
        if (!matchId) return log("[ui] missing matchId");
        currentMatchId = matchId;
        wantReconnect = true;
        reconnectAttempt = 0;
        openWS(matchId, false);
    }

    // Try to restore connection when network/tab comes back.
    window.addEventListener("online", () => {
        if (wantReconnect && currentMatchId && (!ws || ws.readyState === WebSocket.CLOSED)) {
            reconnectAttempt = 0;
            openWS(currentMatchId, true);
        }
    });

    document.addEventListener("visibilitychange", () => {
        if (!document.hidden && wantReconnect && currentMatchId && (!ws || ws.readyState === WebSocket.CLOSED)) {
            reconnectAttempt = 0;
            openWS(currentMatchId, true);
        }
    });
// ---------------- Clipboard
    async function copyText(text) {
        try {
            await navigator.clipboard.writeText(text);
            log("[ui] copied to clipboard");
            return true;
        } catch {
            // fallback
            const ta = document.createElement("textarea");
            ta.value = text;
            document.body.appendChild(ta);
            ta.select();
            document.execCommand("copy");
            ta.remove();
            log("[ui] copied (fallback)");
            return true;
        }
    }

    // ---------------- Auth actions
    $("btnRegister").onclick = async () => {
        try {
            await api("/api/auth/register", {
                method: "POST",
                body: JSON.stringify({
                    email: $("email").value,
                    password: $("password").value,
                    displayName: $("displayName").value
                })
            });
            setAuthMsg(t("authRegistered"), false);
        } catch (e) {
            setAuthMsg(formatError(e), true);
        }
    };

    $("btnLogin").onclick = async () => {
        try {
            await api("/api/auth/login", {
                method: "POST",
                body: JSON.stringify({ email: $("email").value, password: $("password").value })
            });

            await refreshMe();

            setAuthedUI(true);
            setAuthMsg(t("authLoginSuccess"), false);

            // auto-join if match is present (URL > lastMatch)
            const m = getQueryMatchId() || (localStorage.getItem("lastMatchId") || "");
            if (m) {
                setMatchId(m, {persist:true, updateUrl:true});
                connectWS(m);
            }
        } catch (e) {
            setAuthMsg(formatError(e), true);
            setSessionStatus(false);
            setAuthedUI(false);
        }
    };

    async function doLogout() {
        try {
            await api("/api/auth/logout", { method: "POST" });
        } catch {}
        disconnectWS();
        setSessionStatus(false);
        setAuthMsg(t("authLoggedOut"), false);
        setAuthedUI(false);
    }

    $("btnLogout").onclick = () => doLogout();
    $("btnLogoutTop").onclick = () => doLogout();

    $("btnRefreshMe").onclick = async () => {
        try { await refreshMe(); log("[http] refreshed /api/me"); }
        catch (e) { log("[http] refresh /api/me error: " + JSON.stringify(e)); }
    };

    $("btnShowRating").onclick = () => openRating();
    $("btnCloseRating").onclick = () => closeRating();
    $("btnRefreshRating").onclick = () => loadRating();
    $("ratingOverlay").addEventListener("click", (e) => {
        if (e.target && e.target.id === "ratingOverlay") closeRating();
    });

    $("btnRules").onclick = () => openRules();
    $("btnCloseRules").onclick = () => closeRules();
    $("rulesOverlay").addEventListener("click", (e) => {
        if (e.target && e.target.id === "rulesOverlay") closeRules();
    });

    $("btnLang").onclick = () => toggleLang();

    // ---------------- Match actions
    $("matchId").addEventListener("input", () => {
        const m = $("matchId").value.trim();
        setInviteLink(m);
    });
    $("matchId").addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            $("btnJoin").click();
        }
    });
    ["email", "password"].forEach((id) => {
        $(id).addEventListener("keydown", (e) => {
            if (e.key === "Enter") {
                e.preventDefault();
                $("btnLogin").click();
            }
        });
    });

    $("btnCreateAndJoin").onclick = async () => {
        try {
            const out = await api("/api/match", { method: "POST" });
            const m = out.matchId || out.matchID || out.match_id || "";
            if (!m) throw new Error("no matchId in response");
            setMatchId(m, {persist:true, updateUrl:true});
            log("[http] created match: " + m);
            connectWS(m);
        } catch (e) {
            log("[http] create match error: " + JSON.stringify(e));
        }
    };

    $("btnJoin").onclick = async () => {
        const m = $("matchId").value.trim();
        if (!m) return log("[ui] enter matchId");
        setMatchId(m, {persist:true, updateUrl:true});
        connectWS(m);
    };

    $("btnReconnect").onclick = () => {
        const m = localStorage.getItem("lastMatchId") || "";
        if (!m) return log("[ui] no last match");
        setMatchId(m, {persist:true, updateUrl:true});
        connectWS(m);
    };


    $("btnReconnectNow").onclick = () => {
        const m = $("matchId").value.trim() || currentMatchId || (localStorage.getItem("lastMatchId") || "");
        if (!m) return log("[ui] no match to reconnect");
        setMatchId(m, {persist:true, updateUrl:true});
        currentMatchId = m;
        wantReconnect = true;
        reconnectAttempt = 0;
        openWS(m, true);
    };

    $("btnOverlayReconnect").onclick = () => $("btnReconnectNow").click();

    async function copyInviteLink() {
        const link = $("inviteLink").value.trim();
        if (!link) return log("[ui] no invite link");
        await copyText(link);
    }

    $("btnCopyLink").onclick = async () => {
        await copyInviteLink();
    };
    $("btnOverlayCopyLink").onclick = async () => {
        await copyInviteLink();
    };

    $("btnLeave").onclick = () => {
        disconnectWS(true);
        currentMatchId = "";
        lastState = null;

        setMatchId("", {persist:true, updateUrl:true});
        $("youSlot").textContent = "-";
        $("phase").textContent = "-";
        $("round").textContent = "-";
        $("deadline").textContent = "-";
        renderSeries({ p1Wins: 0, p2Wins: 0, draws: 0 }, { p1: "p1", p2: "p2" });
        if ($("matchMode")) $("matchMode").textContent = "-";
        if ($("modeNote")) {
            $("modeNote").textContent = "";
            $("modeNote").classList.add("hidden");
        }
        $("historyBody").innerHTML = `<tr><td colspan="8" class="muted">${t("noData")}</td></tr>`;

        // UX reset
        setStageCopy(t("stageConnectTitle"), t("stageConnectHint"));
        setNetBadge("muted", t("netOffline"));
        setBadge("youRole", "muted", "-");
        setBadge("oppRole", "muted", "-");
        $("youName").textContent = "Вы";
        $("oppName").textContent = "Соперник";
        setBadge("youConn", "muted", "-");
        setBadge("oppConn", "muted", "-");
        setBadge("youSecret", "muted", "-");
        setBadge("oppSecret", "muted", "-");
        setBadge("youGuess", "muted", "-");
        setBadge("oppGuess", "muted", "-");
        $("revealBox").classList.add("hidden");
        setWaitOverlay(false, "", "", false);
        setNetBanner(false);

        $("secret").disabled = true;
        $("btnSetSecret").disabled = true;
        $("guess").disabled = true;
        $("btnGuess").disabled = true;
        $("btnRematch").disabled = true;
        resetActionPanels();

        log("[ui] left match");
    };


// ---------------- Input helpers
    function digits4(v){
        return String(v||"").replace(/\D+/g, "").slice(0, 4);
    }
    ["secret","guess"].forEach((id) => {
        const el = $(id);
        el.addEventListener("input", () => { el.value = digits4(el.value); });
        el.addEventListener("keydown", (e) => {
            if (e.key === "Enter") {
                if (id === "secret" && !$("btnSetSecret").disabled) $("btnSetSecret").click();
                if (id === "guess" && !$("btnGuess").disabled) $("btnGuess").click();
            }
        });
    });

    // ---------------- Game actions
    $("btnSetSecret").onclick = () => {
        const v = digits4($("secret").value);
        $("secret").value = v;
        if (v.length !== 4) return log("[ui] secret must be 4 digits");
        send("set_secret", { secret: v });
        $("secret").value = "";
    };
    $("btnGuess").onclick = () => {
        const v = digits4($("guess").value);
        $("guess").value = v;
        if (v.length !== 4) return log("[ui] guess must be 4 digits");
        send("submit_guess", { guess: v });
        $("guess").value = "";
    };
    $("btnRematch").onclick = () => send("rematch_request", {});

    // ---------------- Boot
    (async function boot() {
        initTheme();
        applyI18n();
        syncAuthPanel(false);
        setSessionStatus(false);

        // Initial UX state
        setNetBadge("muted", t("netOffline"));
        setStageCopy(t("stageConnectTitle"), t("stageConnectHint"));
        $("secret").disabled = true;
        $("btnSetSecret").disabled = true;
        $("guess").disabled = true;
        $("btnGuess").disabled = true;
        $("btnRematch").disabled = true;
        renderSeries({ p1Wins: 0, p2Wins: 0, draws: 0 }, { p1: "p1", p2: "p2" });
        if ($("modeNote")) {
            $("modeNote").textContent = "";
            $("modeNote").classList.add("hidden");
        }
        resetActionPanels();
        const queryMatch = getQueryMatchId();
        if (queryMatch) {
            setMatchId(queryMatch, {persist:false, updateUrl:false});
        }

        try {
            await refreshMe();
            setAuthedUI(true);
            setAuthMsg(t("authTokenRestored"), false);

            // match auto join (URL > lastMatch)
            const m = getQueryMatchId() || (localStorage.getItem("lastMatchId") || "");
            if (m) {
                setMatchId(m, {persist:true, updateUrl:true});
                connectWS(m);
            }
        } catch (e) {
            setSessionStatus(false);
            setAuthedUI(false);
            setAuthMsg(t("authPleaseLogin"), false);
        }
    })();
