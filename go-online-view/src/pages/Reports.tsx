import React, { useContext, useEffect, useMemo, useState } from "react";
import { UserContext } from "../context/UserContext";
import { API_BASE } from "../config";
import Navbar from "../components/Navbar";
import "../styles/Reports.scss";

type ReportRow = {
  id: number;
  game_id: string;
  username: string;
  created_at: string;
};

const Reports: React.FC = () => {
  const { user, loading: userLoading } = useContext(UserContext);

  const [reports, setReports] = useState<ReportRow[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const isStaff = useMemo(
    () => !!user && (user.role === "admin" || user.role === "moderator"),
    [user],
  );

  useEffect(() => {
    if (!isStaff) return;

    const controller = new AbortController();

    (async () => {
      try {
        setLoading(true);
        setError(null);

        const res = await fetch(`${API_BASE}/reports`, {
          method: "GET",
          credentials: "include",
          signal: controller.signal,
          headers: { Accept: "application/json" },
        });

        if (!res.ok) {
          const text = await res.text();
          throw new Error(text || `Failed to load reports (${res.status})`);
        }

        const data: ReportRow[] = await res.json();
        setReports(data);
      } catch (e: unknown) {
        if (e instanceof DOMException && e.name === "AbortError") return;
        setError(e instanceof Error ? e.message : "Something went wrong");
      } finally {
        setLoading(false);
      }
    })();

    return () => controller.abort();
  }, [isStaff]);

  const handleDelete = async (reportId: number) => {
    try {
      setDeletingId(reportId);
      setDeleteError(null);

      const res = await fetch(`${API_BASE}/reports/${reportId}`, {
        method: "DELETE",
        credentials: "include",
        headers: { Accept: "application/json" },
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Failed to delete report (${res.status})`);
      }

      setReports((prev) => prev.filter((r) => r.id !== reportId));
    } catch (e: unknown) {
      setDeleteError(e instanceof Error ? e.message : "Failed to delete report");
    } finally {
      setDeletingId((cur) => (cur === reportId ? null : cur));
    }
  };

  if (userLoading) {
    return (
      <>
        <Navbar />
        <div className="reports-status">Loading…</div>
      </>
    );
  }

  if (!isStaff) {
    return (
      <>
        <Navbar />
        <div className="reports-denied">
          <h1 className="reports-denied__title">Access Denied</h1>
          <p className="reports-denied__text">
            You do not have permission to view this page.
          </p>
        </div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="reports-page">
        <h1 className="reports-title">Reports</h1>

        {loading && <div className="reports-status">Loading reports…</div>}

        {error && (
          <div className="reports-error">
            Error: {error}
          </div>
        )}

        {deleteError && (
          <div className="reports-error">
            Delete error: {deleteError}
          </div>
        )}

        {!loading && !error && reports.length === 0 && (
          <div className="reports-status">No reports found.</div>
        )}

        {!loading && !error && reports.length > 0 && (
          <div className="reports-table">
            <div className="reports-table__header reports-table__grid">
              <div>ID</div>
              <div>Game</div>
              <div>From</div>
              <div>Created</div>
              <div className="reports-table__actionsHeader">Resolved</div>
            </div>

            {reports.map((r) => (
              <div key={r.id} className="reports-table__row reports-table__grid">
                <div className="muted">{r.id}</div>

                <div>
                  {r.game_id ? (
                    <a className="game-link" href={`/game/${r.game_id}`}>
                      {r.game_id}
                    </a>
                  ) : (
                    "—"
                  )}
                </div>

                <div>{r.username ?? "—"}</div>

                <div className="muted">
                  {r.created_at
                    ? new Date(r.created_at)
                        .toISOString()
                        .slice(0, 19)
                        .replace("T", " ")
                    : "—"}
                </div>

                <div className="reports-table__actionsCell">
                  <button
                    type="button"
                    className="delete-btn"
                    onClick={() => handleDelete(r.id)}
                    disabled={deletingId === r.id}
                    aria-label={`Delete report ${r.id}`}
                    title="Delete"
                  >
                    ✓
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
};

export default Reports;
