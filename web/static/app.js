const { useState, useEffect } = React;

/* ---------- API ---------- */
function apiFetch(path, opts = {}) {
    const token = localStorage.getItem("token");
    const headers = opts.headers || {};
    headers["Accept"] = "application/json";
    if (!headers["Content-Type"] && !(opts.body instanceof FormData)) {
        headers["Content-Type"] = "application/json";
    }
    if (token) headers["Authorization"] = "Bearer " + token;

    return fetch(path, { ...opts, headers }).then(async res => {
        const text = await res.text();
        let data = null;
        try { data = text ? JSON.parse(text) : null } catch { data = text }
        if (!res.ok) throw new Error(data?.error || res.status);
        return data;
    });
}

/* ---------- Router ---------- */
function navigate(path) {
    history.pushState({}, "", path);
    window.dispatchEvent(new Event("locationchange"));
}
window.addEventListener("popstate", () =>
    window.dispatchEvent(new Event("locationchange"))
);

/* ---------- Movies ---------- */
function MoviesPage() {
    const [movies, setMovies] = useState([]);

    async function load(title="", year=0) {
        try {
            let list;
            try {
                list = await apiFetch(`/api/movies/search?title=${title}&year=${year}`);
            } catch {
                list = await apiFetch("/api/movies");
                if (title) list = list.filter(m => m.title?.toLowerCase().includes(title.toLowerCase()));
                if (year) list = list.filter(m => m.year === year);
            }
            setMovies(list);
        } catch (e) {
            alert("Movies load failed");
        }
    }

    useEffect(() => { load(); }, []);

    useEffect(() => {
        window.addEventListener("app:search", e => load(e.detail.title, e.detail.year));
    }, []);

    return (
        <div>
            <h2>Movies</h2>
            <div className="movie-grid">
                {movies.map(m => (
                    <div key={m.id} className="card">
                        <h3>
                            <a className="link"
                               href={`/movie/${m.id}`}
                               onClick={e => { e.preventDefault(); navigate(`/movie/${m.id}`); }}>
                                {m.title}
                            </a>
                        </h3>
                        <div className="meta">{m.year} • ⭐ {m.rating || 0}</div>
                        <p>{m.description || "No description"}</p>
                    </div>
                ))}
            </div>
        </div>
    );
}

/* ---------- Movie Detail ---------- */
function MovieDetail({ id }) {
    const [movie, setMovie] = useState(null);
    const [tmdb, setTmdb] = useState(null);

    useEffect(() => {
        async function load() {
            const m = await apiFetch(`/api/movies/${id}`);
            setMovie(m);

            if (m.tmdb_id) {
                try {
                    const t = await apiFetch(`/api/tmdb/movies/${m.tmdb_id}`);
                    setTmdb(t);
                } catch {}
            }
        }
        load();
    }, [id]);

    if (!movie) return <div>Loading...</div>;

    const description =
        movie.description ||
        tmdb?.overview ||
        "No description available";

    let trailerEmbed = null;
    if (tmdb?.trailer_url) {
        const url = tmdb.trailer_url;
        const vid =
            url.includes("v=") ? url.split("v=")[1] :
                url.includes("youtu") ? url.split("/").pop() :
                    null;
        if (vid) trailerEmbed = `https://www.youtube.com/embed/${vid}`;
    }

    return (
        <div>
            <button onClick={() => navigate("/")}>← Back</button>

            <div className="movie-detail">
                <div className="left card">
                    <h2>{movie.title} ({movie.year})</h2>
                    <p>{description}</p>
                    <pre className="json">{JSON.stringify(movie, null, 2)}</pre>
                </div>

                <div className="right card">
                    <h4>Trailer</h4>
                    {trailerEmbed
                        ? <iframe src={trailerEmbed} allowFullScreen />
                        : <div className="small">No trailer</div>
                    }
                    <pre className="json">{JSON.stringify(tmdb, null, 2)}</pre>
                </div>
            </div>
        </div>
    );
}

/* ---------- App ---------- */
function App() {
    const [loc, setLoc] = useState(window.location.pathname);

    useEffect(() => {
        window.addEventListener("locationchange", () =>
            setLoc(window.location.pathname)
        );
    }, []);

    let page = <MoviesPage />;
    if (loc.startsWith("/movie/")) {
        page = <MovieDetail id={loc.split("/")[2]} />;
    }

    return page;
}

ReactDOM.createRoot(document.getElementById("root")).render(<App />);
