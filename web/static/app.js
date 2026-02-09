/* MovieHub frontend (vanilla JS) */

const API = {
    movies: "/api/movies",
    search: "/api/movies/search",
    movie: (id) => `/api/movies/${id}`,
    tmdbMovie: (tmdbId) => `/api/tmdb/movies/${tmdbId}`,
    reviews: (movieId) => `/api/movies/${movieId}/reviews`,
    addReview: (movieId) => `/api/movies/${movieId}/reviews`,
    updateReview: (movieId) => `/api/movies/${movieId}/reviews`,
    deleteReview: (movieId) => `/api/movies/${movieId}/reviews`,
    adminDeleteReview: (reviewId) => `/api/reviews/${reviewId}`,

    authRegister: "/api/auth/register",
    authLogin: "/api/auth/login",

    me: "/api/me",
    mePassword: "/api/me/password",
    meDelete: "/api/me",
    userById: (id) => `/api/users/${id}`,
};

const state = {
    token: localStorage.getItem("token") || "",
    me: null,
    movies: [],
    selectedMovie: null,
    selectedReviews: [],
    userCache: new Map(),
};

const $ = (id) => document.getElementById(id);

function toast(msg) {
    const el = $("toast");
    el.textContent = msg;
    el.classList.remove("hidden");
    clearTimeout(toast._t);
    toast._t = setTimeout(() => el.classList.add("hidden"), 2400);
}

function openModal(id) { $(id).classList.remove("hidden"); }
function closeModal(id) { $(id).classList.add("hidden"); }

async function apiFetch(path, { method="GET", body=null, auth=false } = {}) {
    const headers = { "Accept": "application/json" };
    if (body !== null) headers["Content-Type"] = "application/json";
    if (auth && state.token) headers["Authorization"] = `Bearer ${state.token}`;

    const res = await fetch(path, {
        method,
        headers,
        body: body === null ? null : JSON.stringify(body),
    });

    let data = null;
    const ct = res.headers.get("content-type") || "";
    if (ct.includes("application/json")) {
        try { data = await res.json(); } catch { data = null; }
    }

    if (!res.ok) {
        const errMsg = (data && (data.error || data.message))
            ? (data.error || data.message)
            : `${res.status} ${res.statusText}`;
        throw new Error(errMsg);
    }

    return data;
}

/* ---------- YouTube embed helper ---------- */

/* ---------- YouTube embed helper (improved) ---------- */

function toYouTubeEmbed(data) {
    // Accepts either:
    // - a full URL string ("https://www.youtube.com/watch?v=xxx")
    // - a TMDB object { trailer_url: "...", key: "...", Results: [...] }
    // - a TMDB videos response { results: [...] }
    // Returns embed URL string or null.

    // if data is falsy
    if (!data) return null;

    // If it's an object with trailer_url
    if (typeof data === "object") {
        if (data.trailer_url && typeof data.trailer_url === "string") {
            // fallthrough to parse string below
            data = data.trailer_url;
        } else if (data.key && typeof data.key === "string") {
            return `https://www.youtube.com/embed/${data.key}`;
        } else if (Array.isArray(data.results) && data.results.length > 0) {
            // find first YouTube Trailer or Teaser
            const v = data.results.find(x => (
                (x.site && x.site.toLowerCase().includes("youtube")) &&
                (x.type && (x.type.toLowerCase().includes("trailer") || x.type.toLowerCase().includes("teaser")))
            )) || data.results[0];
            if (v && v.key) return `https://www.youtube.com/embed/${v.key}`;
        }
        // no known shape
        return null;
    }

    // At this point `data` should be a string (url or key)
    if (typeof data !== "string") return null;

    // try to parse as URL
    try {
        const u = new URL(data);
        // short link youtu.be/<id>
        if (u.hostname.includes("youtu.be")) {
            const id = u.pathname.replace("/", "");
            return id ? `https://www.youtube.com/embed/${id}` : null;
        }
        // regular youtube link youtube.com/watch?v=<id>
        if (u.hostname.includes("youtube.com") || u.hostname.includes("www.youtube.com")) {
            const id = u.searchParams.get("v");
            return id ? `https://www.youtube.com/embed/${id}` : null;
        }
        // sometimes TMDB might give an embed path already, return as-is if it looks like embed
        if (u.hostname.includes("youtube-nocookie.com") || u.pathname.includes("/embed/")) {
            return data;
        }
    } catch (e) {
        // not a URL — maybe plain id
        const plain = data.trim();
        if (plain.length >= 6 && plain.length <= 20 && /^[A-Za-z0-9_\-]+$/.test(plain)) {
            return `https://www.youtube.com/embed/${plain}`;
        }
        return null;
    }
    return null;
}


/* ---------- Auth ---------- */

function setToken(token) {
    state.token = token || "";
    if (state.token) localStorage.setItem("token", state.token);
    else localStorage.removeItem("token");
}

async function loadMe() {
    if (!state.token) {
        state.me = null;
        return;
    }
    try {
        state.me = await apiFetch(API.me, { auth: true });
    } catch {
        setToken("");
        state.me = null;
    }
}

function renderAuthUI() {
    const authArea = $("authArea");
    const userArea = $("userArea");
    const userBadge = $("userBadge");

    if (!state.token || !state.me) {
        authArea.classList.remove("hidden");
        userArea.classList.add("hidden");
        return;
    }

    authArea.classList.add("hidden");
    userArea.classList.remove("hidden");
    userBadge.textContent = `${state.me.username} (${state.me.role})`;
}

async function login(email, password) {
    const out = await apiFetch(API.authLogin, { method:"POST", body:{ email, password } });
    setToken(out.token);
    await loadMe();
    renderAuthUI();
    updateComposerVisibility();
    toast("Logged in ✅");
}

async function register(username, email, password) {
    await apiFetch(API.authRegister, { method:"POST", body:{ username, email, password } });
    toast("Registered ✅ Now log in.");
}

function logout() {
    setToken("");
    state.me = null;
    state.userCache.clear();
    renderAuthUI();
    updateComposerVisibility();
    toast("Logged out");
}

/* ---------- Movies ---------- */

function movieCard(m) {
    const el = document.createElement("div");
    el.className = "item";
    el.dataset.id = String(m.id);

    const title = document.createElement("div");
    title.className = "item-title";
    title.textContent = m.title || "Untitled";

    const meta = document.createElement("div");
    meta.className = "item-meta";
    meta.innerHTML = `
    <span>${m.year ?? "—"}</span>
    <span>★ ${Number(m.rating || 0).toFixed(1)}</span>
    <span class="mono">id:${m.id}</span>
  `;

    el.appendChild(title);
    el.appendChild(meta);

    el.addEventListener("click", () => selectMovie(m.id));
    return el;
}

function renderMovieList() {
    const list = $("movieList");
    list.innerHTML = "";

    $("movieCount").textContent = `${state.movies.length} items`;

    for (const m of state.movies) {
        const el = movieCard(m);
        if (state.selectedMovie && state.selectedMovie.id === m.id) el.classList.add("active");
        list.appendChild(el);
    }
}

async function loadMovies() {
    state.movies = await apiFetch(API.movies);
    renderMovieList();
}

async function searchMovies(title, year) {
    const qs = new URLSearchParams();
    if (title) qs.set("title", title);
    if (year) qs.set("year", String(year));
    state.movies = await apiFetch(`${API.search}?${qs.toString()}`);
    renderMovieList();
}

function renderMovieView(movie, trailerData=null) {
    $("emptyState").classList.add("hidden");
    $("movieView").classList.remove("hidden");

    $("mvTitle").textContent = movie.title || "—";
    $("mvYear").textContent = movie.year ? String(movie.year) : "—";
    $("mvRating").textContent = `★ ${Number(movie.rating || 0).toFixed(1)}`;
    $("mvTMDB").textContent = `TMDB: ${movie.tmdb_id || 0}`;
    $("mvDesc").textContent = movie.description || "No description.";

    // Trailer (robust handling)
    const trailerBox = $("trailerBox");
    const trailerLink = $("trailerLink");
    const trailerFrame = $("trailerFrame");

    // Clear previous src to stop playback when switching movies
    trailerFrame.src = "";

    // Compute embed URL using the helper (supports multiple TMDB shapes)
    let embedUrl = null;
    let publicUrl = null;

    if (trailerData) {
        // trailerData might be: { trailer_url: "...", key: "..." } or { results: [...] } or plain string
        if (typeof trailerData === "string") {
            publicUrl = trailerData;
            embedUrl = toYouTubeEmbed(trailerData);
        } else if (typeof trailerData === "object") {
            // If it has a top-level trailer_url, use that as public url
            if (trailerData.trailer_url) publicUrl = trailerData.trailer_url;
            // If it has key property
            if (trailerData.key) embedUrl = toYouTubeEmbed(trailerData.key);
            // If it has results (TMDB videos response)
            if (!embedUrl && Array.isArray(trailerData.results)) embedUrl = toYouTubeEmbed(trailerData);
            // fallback: try full object (some endpoints return { trailer_url: .. })
            if (!embedUrl && publicUrl) embedUrl = toYouTubeEmbed(publicUrl);
        }
    }

    // If we don't have embedUrl but have a publicUrl, still set the link (user can open manually)
    if (publicUrl) {
        trailerLink.href = publicUrl;
        trailerLink.classList.remove("hidden");
    } else {
        trailerLink.href = "#";
        trailerLink.classList.add("hidden");
    }

    if (embedUrl) {
        // add autoplay=0 by default; if you want autoplay set to 1, append ?autoplay=1 (careful with UX)
        // ensure we don't duplicate querystring
        const sep = embedUrl.includes("?") ? "&" : "?";
        trailerFrame.src = embedUrl + sep + "rel=0"; // rel=0 to reduce related videos
        trailerBox.classList.remove("hidden");
    } else {
        trailerFrame.removeAttribute("src");
        trailerBox.classList.add("hidden");
    }


    renderMovieAdminActions(movie);
}

function renderMovieAdminActions(movie) {
    const box = $("movieAdminActions");
    box.innerHTML = "";

    // Keep edit/delete visible only if logged in (your backend should also enforce admin-only)
    if (!state.me) return;

    const btnEdit = document.createElement("button");
    btnEdit.className = "btn btn-ghost";
    btnEdit.textContent = "Quick Edit";
    btnEdit.onclick = async () => {
        const title = prompt("New title (leave empty to keep):", movie.title || "");
        const yearStr = prompt("New year (leave empty to keep):", movie.year ? String(movie.year) : "");
        const description = prompt("New description (leave empty to keep):", movie.description || "");

        const body = {};
        if (title && title.trim()) body.title = title.trim();
        const year = parseInt(yearStr, 10);
        if (!Number.isNaN(year) && year > 0) body.year = year;
        if (description && description.trim()) body.description = description.trim();

        if (Object.keys(body).length === 0) return;

        try {
            const updated = await apiFetch(API.movie(movie.id), { method:"PUT", body, auth:true });
            toast("Movie updated ✅");
            await loadMovies();
            await selectMovie(updated.id);
        } catch (e) {
            toast(`Update failed: ${e.message}`);
        }
    };

    const btnDelete = document.createElement("button");
    btnDelete.className = "btn btn-danger";
    btnDelete.textContent = "Delete movie";
    btnDelete.onclick = async () => {
        if (!confirm("Delete this movie?")) return;
        try {
            await apiFetch(API.movie(movie.id), { method:"DELETE", auth:true });
            toast("Movie deleted ✅");
            state.selectedMovie = null;
            await loadMovies();
            $("movieView").classList.add("hidden");
            $("emptyState").classList.remove("hidden");
        } catch (e) {
            toast(`Delete failed: ${e.message}`);
        }
    };

    box.appendChild(btnEdit);
    box.appendChild(btnDelete);
}

async function selectMovie(id) {
    try {
        // 1) Load movie details
        const movie = await apiFetch(API.movie(id));
        state.selectedMovie = movie;
        renderMovieList();

        // 2) Render immediately (do NOT depend on TMDB)
        renderMovieView(movie, null);

        // 3) Load reviews
        await loadReviews(movie.id);
        updateComposerVisibility();

        // 4) Trailer fetch in background (ignore failure)
        if (movie.tmdb_id && movie.tmdb_id > 0) {
            apiFetch(API.tmdbMovie(movie.tmdb_id))
                .then((trailer) => renderMovieView(movie, trailer))
                .catch(() => {/* ignore */});
        }
    } catch (e) {
        toast(`Failed to load movie: ${e.message}`);
    }
}

/* ---------- Reviews ---------- */

async function loadReviews(movieId) {
    try {
        state.selectedReviews = await apiFetch(API.reviews(movieId));
    } catch (e) {
        // If backend ever returns 404, treat as empty list (but we fixed backend handler too)
        state.selectedReviews = [];
    }

    $("reviewCount").textContent = `${state.selectedReviews.length} reviews`;
    await renderReviews();
}

function reviewActions(review) {
    const row = document.createElement("div");
    row.className = "row";

    if (state.me && state.me.role === "admin") {
        const btnAdminDel = document.createElement("button");
        btnAdminDel.className = "btn btn-danger";
        btnAdminDel.textContent = `Admin delete #${review.id}`;
        btnAdminDel.onclick = async () => {
            if (!confirm(`Admin delete review #${review.id}?`)) return;
            try {
                await apiFetch(API.adminDeleteReview(review.id), { method:"DELETE", auth:true });
                toast("Review deleted ✅");
                await loadReviews(state.selectedMovie.id);
            } catch (e) {
                toast(`Admin delete failed: ${e.message}`);
            }
        };
        row.appendChild(btnAdminDel);
    }

    return row;
}

async function ensureUserLabel(userId) {
    if (!state.me) return `user:${userId}`;
    if (state.userCache.has(userId)) return state.userCache.get(userId);

    try {
        const u = await apiFetch(API.userById(userId), { auth:true });
        const label = u.username ? `${u.username} (id:${u.id})` : `user:${userId}`;
        state.userCache.set(userId, label);
        return label;
    } catch {
        return `user:${userId}`;
    }
}

async function renderReviews() {
    const list = $("reviewList");
    list.innerHTML = "";

    if (!state.selectedReviews.length) {
        const empty = document.createElement("div");
        empty.className = "muted";
        empty.textContent = "No reviews yet.";
        list.appendChild(empty);
        return;
    }

    for (const r of state.selectedReviews) {
        const el = document.createElement("div");
        el.className = "review";

        const top = document.createElement("div");
        top.className = "review-top";

        const user = document.createElement("div");
        user.className = "review-user";
        user.textContent = await ensureUserLabel(r.userId);

        const score = document.createElement("div");
        score.className = "review-score";
        score.textContent = `Score: ${r.score}`;

        top.appendChild(user);
        top.appendChild(score);

        const text = document.createElement("p");
        text.className = "review-text";
        text.textContent = r.text || "";

        el.appendChild(top);
        if (r.text) el.appendChild(text);
        el.appendChild(reviewActions(r));

        list.appendChild(el);
    }
}

function updateComposerVisibility() {
    const composer = $("reviewComposer");
    const hint = $("reviewHint");

    if (!state.me || !state.token) {
        composer.classList.add("hidden");
        hint.textContent = "Log in to post a review.";
        return;
    }
    composer.classList.remove("hidden");
    hint.textContent = `Posting as ${state.me.username} (${state.me.role})`;
}

async function addReview() {
    if (!state.selectedMovie) return;

    const score = parseInt($("revScore").value, 10);
    const text = ($("revText").value || "").trim();

    // basic validation before sending
    if (Number.isNaN(score) || score < 1 || score > 10) {
        toast("Score must be between 1 and 10");
        return;
    }

    try {
        await apiFetch(API.addReview(state.selectedMovie.id), {
            method: "POST",
            // send BOTH keys so backend can bind either json:"text" or json:"comment"
            body: { score, text, comment: text },
            auth: true
        });

        toast("Review added ✅");
        $("revText").value = "";
        await loadReviews(state.selectedMovie.id);
        await loadMovies();
    } catch (e) {
        toast(`Add review failed: ${e.message}`);
    }
}

async function updateMyScoreOnly() {
    if (!state.selectedMovie) return;
    const score = parseInt($("revScore").value, 10);
    try {
        await apiFetch(API.updateReview(state.selectedMovie.id), {
            method:"PUT",
            body:{ score },
            auth:true
        });
        toast("Score updated ✅");
        await loadReviews(state.selectedMovie.id);
        await loadMovies();
    } catch (e) {
        toast(`Update failed: ${e.message}`);
    }
}

async function deleteMyReview() {
    if (!state.selectedMovie) return;
    if (!confirm("Delete your review for this movie?")) return;
    try {
        await apiFetch(API.deleteReview(state.selectedMovie.id), {
            method:"DELETE",
            auth:true
        });
        toast("Review deleted ✅");
        await loadReviews(state.selectedMovie.id);
        await loadMovies();
    } catch (e) {
        toast(`Delete failed: ${e.message}`);
    }
}

/* ---------- Profile ---------- */

async function openProfile() {
    try {
        const me = await apiFetch(API.me, { auth:true });
        state.me = me;
        $("meSummary").textContent = `${me.username} <${me.email}>`;
        $("meRole").textContent = me.role || "user";
        $("meUsername").value = me.username || "";
        $("meEmail").value = me.email || "";
        openModal("profileModal");
    } catch (e) {
        toast(`Profile load failed: ${e.message}`);
    }
}

async function saveProfile() {
    const username = $("meUsername").value.trim();
    const email = $("meEmail").value.trim();
    try {
        const updated = await apiFetch(API.me, {
            method:"PUT",
            body:{ username, email },
            auth:true
        });
        state.me = updated;
        renderAuthUI();
        toast("Profile updated ✅");
    } catch (e) {
        toast(`Update failed: ${e.message}`);
    }
}

async function changePassword() {
    const password = $("mePass").value;
    try {
        await apiFetch(API.mePassword, {
            method:"PUT",
            body:{ password },
            auth:true
        });
        $("mePass").value = "";
        toast("Password changed ✅");
    } catch (e) {
        toast(`Password change failed: ${e.message}`);
    }
}

async function deleteAccount() {
    if (!confirm("Delete your account permanently?")) return;
    try {
        await apiFetch(API.meDelete, { method:"DELETE", auth:true });
        logout();
        closeModal("profileModal");
        toast("Account deleted");
    } catch (e) {
        toast(`Delete failed: ${e.message}`);
    }
}

/* ---------- Wire up UI ---------- */

function wireUI() {
    $("btnSearch").addEventListener("click", async () => {
        const title = $("qTitle").value.trim();
        const year = parseInt($("qYear").value, 10);
        try {
            await searchMovies(title, Number.isNaN(year) ? 0 : year);
            toast("Search done");
        } catch (e) {
            toast(`Search failed: ${e.message}`);
        }
    });

    $("btnClear").addEventListener("click", async () => {
        $("qTitle").value = "";
        $("qYear").value = "";
        await loadMovies();
    });

    $("btnLoginOpen").addEventListener("click", () => openModal("loginModal"));
    $("btnRegisterOpen").addEventListener("click", () => openModal("registerModal"));

    document.querySelectorAll("[data-close]").forEach(btn => {
        btn.addEventListener("click", () => closeModal(btn.dataset.close));
    });

    document.querySelectorAll(".modal").forEach(m => {
        m.addEventListener("click", (e) => {
            if (e.target === m) m.classList.add("hidden");
        });
    });

    $("btnLogin").addEventListener("click", async () => {
        const email = $("loginEmail").value.trim();
        const password = $("loginPass").value;
        try {
            await login(email, password);
            closeModal("loginModal");
        } catch (e) {
            toast(`Login failed: ${e.message}`);
        }
    });

    $("btnRegister").addEventListener("click", async () => {
        const username = $("regUsername").value.trim();
        const email = $("regEmail").value.trim();
        const password = $("regPass").value;
        try {
            await register(username, email, password);
            closeModal("registerModal");
        } catch (e) {
            toast(`Register failed: ${e.message}`);
        }
    });

    $("btnLogout").addEventListener("click", logout);
    $("btnProfile").addEventListener("click", openProfile);

    $("btnUpdateProfile").addEventListener("click", saveProfile);
    $("btnChangePassword").addEventListener("click", changePassword);
    $("btnDeleteAccount").addEventListener("click", deleteAccount);

    $("btnAddReview").addEventListener("click", addReview);
    $("btnUpdateScore").addEventListener("click", updateMyScoreOnly);
    $("btnDeleteMyReview").addEventListener("click", deleteMyReview);

    $("qTitle").addEventListener("keydown", (e) => { if (e.key === "Enter") $("btnSearch").click(); });
    $("qYear").addEventListener("keydown", (e) => { if (e.key === "Enter") $("btnSearch").click(); });
}

async function boot() {
    wireUI();
    await loadMe();
    renderAuthUI();
    updateComposerVisibility();

    try {
        await loadMovies();
    } catch (e) {
        toast(`Failed to load movies: ${e.message}`);
    }
}

boot();
