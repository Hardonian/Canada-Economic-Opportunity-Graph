"use client";

import Link from "next/link";
import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { ArrowUpRight, FileText, FolderLock, Loader2, MessageSquareText, Plus, ShieldCheck, Upload, UserPlus, UsersRound } from "lucide-react";
import type { Project } from "@/lib/types";
import type { DealRoomActor } from "@/lib/deal-room-gateway";

type Membership = { user_id: string; role: "OWNER" | "ADMIN" | "MEMBER" | "VIEWER"; added_by: string; created_at: string };
type Room = { id: string; project_id: string; name: string; classification: string; retention_days: number; created_by: string; created_at: string; updated_at: string; members: Record<string, Membership> };
type Document = { id: string; title: string; file_name: string; content_type: string; size_bytes: number; sha256: string; classification: string; uploaded_by: string; created_at: string; retention_at: string };
type Comment = { id: string; document_id?: string; body: string; created_by: string; created_at: string };
type AuditEvent = { id: string; action: string; subject_id: string; actor_id: string; occurred_at: string; metadata?: Record<string, unknown> };

const CLASSIFICATIONS = ["UNCLASSIFIED", "PROTECTED_A", "PROTECTED_B", "PROTECTED_C", "CONFIDENTIAL"];

function roomRequest(path: string, init?: RequestInit) {
  return fetch(`/api/deal-room/${path}`, { cache: "no-store", ...init, headers: { Accept: "application/json", ...(init?.headers || {}) } });
}

async function json<T>(response: Response): Promise<T> {
  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    const detail = payload?.error?.message || payload?.error || "Private deal-room request failed.";
    throw new Error(detail);
  }
  return payload as T;
}

function bytes(value: number) {
  if (value < 1_000_000) return `${Math.max(1, Math.round(value / 1024))} KB`;
  return `${(value / 1_000_000).toFixed(1)} MB`;
}

export default function DealRoomClient({ actor, projects }: { actor: DealRoomActor; projects: Project[] }) {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [selectedRoomID, setSelectedRoomID] = useState("");
  const [documents, setDocuments] = useState<Document[]>([]);
  const [comments, setComments] = useState<Comment[]>([]);
  const [audit, setAudit] = useState<AuditEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [roomName, setRoomName] = useState("");
  const [projectID, setProjectID] = useState(projects[0]?.id ?? "");
  const [classification, setClassification] = useState("PROTECTED_A");
  const [retention, setRetention] = useState(2555);
  const [memberID, setMemberID] = useState("");
  const [memberRole, setMemberRole] = useState("MEMBER");
  const [commentBody, setCommentBody] = useState("");

  const selectedRoom = useMemo(() => rooms.find((room) => room.id === selectedRoomID) ?? rooms[0], [rooms, selectedRoomID]);
  const ownMembership = selectedRoom?.members[actor.user_id];
  const canManage = ownMembership?.role === "OWNER" || ownMembership?.role === "ADMIN";
  const canWrite = canManage || ownMembership?.role === "MEMBER";

  const loadRooms = useCallback(async () => {
    const data = await json<{ rooms: Room[] }>(await roomRequest("rooms"));
    setRooms(data.rooms);
    setSelectedRoomID((current) => current || data.rooms[0]?.id || "");
  }, []);

  const loadRoomContext = useCallback(async (roomID: string) => {
    if (!roomID) { setDocuments([]); setComments([]); setAudit([]); return; }
    const [docs, roomComments, roomAudit] = await Promise.all([
      json<{ documents: Document[] }>(await roomRequest(`rooms/${roomID}/documents`)),
      json<{ comments: Comment[] }>(await roomRequest(`rooms/${roomID}/comments`)),
      json<{ audit: AuditEvent[] }>(await roomRequest(`rooms/${roomID}/audit`)),
    ]);
    setDocuments(docs.documents); setComments(roomComments.comments); setAudit(roomAudit.audit);
  }, []);

  useEffect(() => {
    loadRooms().catch((error) => setMessage(error instanceof Error ? error.message : "Unable to load private rooms.")).finally(() => setLoading(false));
  }, [loadRooms]);

  useEffect(() => {
    if (!selectedRoomID) return;
    loadRoomContext(selectedRoomID).catch((error) => setMessage(error instanceof Error ? error.message : "Unable to load room context."));
  }, [loadRoomContext, selectedRoomID]);

  async function createRoom(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setSaving(true); setMessage("");
    try {
      const data = await json<{ room: Room }>(await roomRequest("rooms", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ project_id: projectID, name: roomName, classification, retention_days: retention }) }));
      setRooms((items) => [data.room, ...items]); setSelectedRoomID(data.room.id); setRoomName(""); setMessage("Encrypted deal room created. You are its owner.");
    } catch (error) { setMessage(error instanceof Error ? error.message : "Unable to create deal room."); } finally { setSaving(false); }
  }

  async function addMember(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); if (!selectedRoom) return; setSaving(true); setMessage("");
    try {
      await json(await roomRequest(`rooms/${selectedRoom.id}/members`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ user_id: memberID, role: memberRole }) }));
      await loadRooms(); setMemberID(""); setMessage("Member granted access. Their identity must carry the same tenant claim.");
    } catch (error) { setMessage(error instanceof Error ? error.message : "Unable to grant member access."); } finally { setSaving(false); }
  }

  async function uploadDocument(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); if (!selectedRoom) return; const form = event.currentTarget; const file = (new FormData(form).get("file") as File | null); if (!file) { setMessage("Select a PDF, Office document, text, or CSV file."); return; }
    setSaving(true); setMessage("");
    try {
      const data = new FormData(form); data.set("classification", selectedRoom.classification);
      await json(await roomRequest(`rooms/${selectedRoom.id}/documents`, { method: "POST", body: data }));
      form.reset(); await loadRoomContext(selectedRoom.id); setMessage("Document encrypted at rest and recorded in the room audit trail.");
    } catch (error) { setMessage(error instanceof Error ? error.message : "Unable to upload document."); } finally { setSaving(false); }
  }

  async function addComment(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); if (!selectedRoom || commentBody.trim().length < 2) return; setSaving(true); setMessage("");
    try {
      await json(await roomRequest(`rooms/${selectedRoom.id}/comments`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ body: commentBody }) }));
      setCommentBody(""); await loadRoomContext(selectedRoom.id);
    } catch (error) { setMessage(error instanceof Error ? error.message : "Unable to post comment."); } finally { setSaving(false); }
  }

  return <main className="mx-auto max-w-7xl space-y-6 px-4 py-8 sm:px-6 lg:px-8">
    <section className="relative overflow-hidden rounded-3xl border border-primary/35 bg-gradient-to-br from-card to-[#07130d] p-7 shadow-2xl shadow-black/40"><div className="pointer-events-none absolute -right-24 -top-24 h-72 w-72 rounded-full bg-aurora/15 blur-3xl" />
      <div className="relative flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between"><div><p className="font-mono text-xs font-bold tracking-[0.18em] text-aurora">PRIVATE · TENANT-ISOLATED · AUDITABLE</p><h1 className="mt-3 text-3xl font-black tracking-tight text-text-main sm:text-5xl">Deal Room</h1><p className="mt-3 max-w-2xl text-sm leading-relaxed text-text-muted">Private diligence, evidence exchange, and accountable collaboration for major projects. This data never enters the public opportunity graph.</p></div><div className="rounded-2xl border border-border bg-background/65 px-4 py-3 text-xs"><p className="font-mono text-text-subtle">TENANT</p><p className="mt-1 font-bold text-text-main">{actor.organization}</p><p className="mt-1 font-mono text-aurora">{actor.clearance} · {actor.user_id}</p></div></div>
    </section>

    <div className="grid gap-6 xl:grid-cols-[0.75fr_1.25fr]">
      <aside className="space-y-6"><section className="glass-panel rounded-2xl border border-border p-5"><div className="flex items-center gap-2"><Plus className="h-4 w-4 text-aurora" /><h2 className="text-sm font-bold text-text-main">Create private room</h2></div><form onSubmit={createRoom} className="mt-4 space-y-3"><input required value={roomName} onChange={(event) => setRoomName(event.target.value)} maxLength={160} placeholder="e.g. Project financing diligence" className="w-full rounded-xl border border-border bg-background px-3 py-2.5 text-xs text-text-main placeholder:text-text-subtle" /><select value={projectID} onChange={(event) => setProjectID(event.target.value)} className="w-full rounded-xl border border-border bg-background px-3 py-2.5 text-xs text-text-main">{projects.map((project) => <option key={project.id} value={project.id}>{project.name}</option>)}</select><div className="grid grid-cols-2 gap-2"><select value={classification} onChange={(event) => setClassification(event.target.value)} className="rounded-xl border border-border bg-background px-3 py-2.5 text-xs text-text-main">{CLASSIFICATIONS.map((level) => <option key={level} value={level}>{level}</option>)}</select><select value={retention} onChange={(event) => setRetention(Number(event.target.value))} className="rounded-xl border border-border bg-background px-3 py-2.5 text-xs text-text-main"><option value={365}>1 year</option><option value={1095}>3 years</option><option value={2555}>7 years</option></select></div><button disabled={saving} className="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-primary px-4 py-2.5 text-xs font-black text-[#06110c] disabled:opacity-60"><FolderLock className="h-3.5 w-3.5" />Create encrypted room</button></form></section>
        <section className="glass-panel rounded-2xl border border-border p-3"><div className="flex items-center justify-between px-2 py-2"><h2 className="text-sm font-bold text-text-main">Your rooms</h2><span className="font-mono text-xs text-aurora">{rooms.length}</span></div><div className="space-y-1">{rooms.map((room) => <button key={room.id} onClick={() => setSelectedRoomID(room.id)} className={`w-full rounded-xl border p-3 text-left transition ${room.id === selectedRoom?.id ? "border-primary/55 bg-primary/10" : "border-transparent hover:border-borderSubtle hover:bg-surface"}`}><p className="truncate text-xs font-bold text-text-main">{room.name}</p><p className="mt-1 font-mono text-[10px] text-text-subtle">{room.classification} · {Object.keys(room.members).length} members</p></button>)}{!loading && !rooms.length && <p className="p-4 text-xs text-text-muted">Create the first private diligence room for this tenant.</p>}</div></section></aside>
      <section className="space-y-6"><section className="glass-panel rounded-2xl border border-border p-5 sm:p-6">{loading ? <div className="flex min-h-40 items-center justify-center text-text-muted"><Loader2 className="mr-2 h-4 w-4 animate-spin" />Loading private rooms…</div> : selectedRoom ? <><div className="flex flex-col gap-3 border-b border-borderSubtle pb-4 sm:flex-row sm:items-start sm:justify-between"><div><p className="font-mono text-[10px] text-aurora">{selectedRoom.classification} · {selectedRoom.retention_days}-DAY RETENTION</p><h2 className="mt-1 text-xl font-black text-text-main">{selectedRoom.name}</h2><p className="mt-1 text-xs text-text-muted">Project {selectedRoom.project_id} · all document actions are audit logged.</p></div><Link href={`/projects/${projects.find((project) => project.id === selectedRoom.project_id)?.slug ?? ""}`} className="inline-flex items-center gap-1 text-xs font-bold text-aurora hover:underline">Public project intelligence <ArrowUpRight className="h-3.5 w-3.5" /></Link></div>
          <div className="mt-5 grid gap-5 lg:grid-cols-2"><div><div className="flex items-center gap-2"><FileText className="h-4 w-4 text-gold" /><h3 className="text-sm font-bold text-text-main">Encrypted documents</h3></div><div className="mt-3 space-y-2">{documents.map((document) => <div key={document.id} className="rounded-xl border border-borderSubtle bg-surface/55 p-3"><div className="flex items-start justify-between gap-3"><div className="min-w-0"><p className="truncate text-xs font-bold text-text-main">{document.title}</p><p className="mt-1 font-mono text-[10px] text-text-subtle">{document.classification} · {bytes(document.size_bytes)} · SHA-256 {document.sha256.slice(0, 12)}…</p></div><a href={`/api/deal-room/rooms/${selectedRoom.id}/documents/${document.id}`} className="rounded-lg border border-primary/40 px-2 py-1 text-[10px] font-bold text-aurora">Download</a></div></div>)}{!documents.length && <p className="rounded-xl border border-dashed border-borderSubtle p-4 text-xs text-text-muted">No private documents yet.</p>}</div>{canWrite && <form onSubmit={uploadDocument} className="mt-3 rounded-xl border border-borderSubtle bg-background/55 p-3"><input name="title" required maxLength={180} placeholder="Document title" className="w-full rounded-lg border border-border bg-background px-2.5 py-2 text-xs text-text-main" /><input name="file" required type="file" accept=".pdf,.txt,.csv,.docx,.xlsx,application/pdf,text/plain,text/csv,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" className="mt-2 block w-full text-[11px] text-text-muted file:mr-3 file:rounded-lg file:border-0 file:bg-primary/15 file:px-2 file:py-1.5 file:text-xs file:font-bold file:text-aurora" /><button disabled={saving} className="mt-3 inline-flex items-center gap-2 text-xs font-bold text-aurora"><Upload className="h-3.5 w-3.5" />Upload encrypted document</button></form>}</div>
            <div><div className="flex items-center gap-2"><MessageSquareText className="h-4 w-4 text-sky-300" /><h3 className="text-sm font-bold text-text-main">Collaboration notes</h3></div><div className="mt-3 max-h-72 space-y-2 overflow-y-auto">{comments.map((comment) => <article key={comment.id} className="rounded-xl border border-borderSubtle bg-surface/55 p-3"><p className="text-xs text-text-main">{comment.body}</p><p className="mt-2 font-mono text-[10px] text-text-subtle">{comment.created_by} · {new Date(comment.created_at).toLocaleString("en-CA")}</p></article>)}{!comments.length && <p className="rounded-xl border border-dashed border-borderSubtle p-4 text-xs text-text-muted">No collaboration notes yet.</p>}</div>{canWrite && <form onSubmit={addComment} className="mt-3 flex gap-2"><input value={commentBody} onChange={(event) => setCommentBody(event.target.value)} maxLength={2000} placeholder="Record a diligence note or accountable next action" className="min-w-0 flex-1 rounded-xl border border-border bg-background px-3 py-2 text-xs text-text-main" /><button disabled={saving} className="rounded-xl border border-primary/45 bg-primary/10 px-3 text-xs font-bold text-aurora">Post</button></form>}</div></div>
          {canManage && <form onSubmit={addMember} className="mt-5 flex flex-col gap-2 border-t border-borderSubtle pt-5 sm:flex-row"><div className="flex items-center gap-2 text-xs font-bold text-text-main"><UserPlus className="h-4 w-4 text-aurora" />Grant room membership</div><input required value={memberID} onChange={(event) => setMemberID(event.target.value)} placeholder="Authenticated user ID" className="min-w-0 flex-1 rounded-xl border border-border bg-background px-3 py-2 text-xs text-text-main" /><select value={memberRole} onChange={(event) => setMemberRole(event.target.value)} className="rounded-xl border border-border bg-background px-3 py-2 text-xs text-text-main"><option>MEMBER</option><option>VIEWER</option><option>ADMIN</option></select><button disabled={saving} className="rounded-xl border border-primary/45 bg-primary/10 px-3 py-2 text-xs font-bold text-aurora">Grant</button></form>}</> : <div className="flex min-h-64 flex-col items-center justify-center text-center"><UsersRound className="h-8 w-8 text-text-subtle" /><h2 className="mt-3 font-bold text-text-main">No room selected</h2><p className="mt-1 max-w-sm text-xs text-text-muted">Create a private room to start tenant-isolated diligence and collaboration.</p></div>}</section>
        {selectedRoom && <section className="rounded-2xl border border-border bg-surface/50 p-5"><div className="flex items-center gap-2"><ShieldCheck className="h-4 w-4 text-aurora" /><h2 className="text-sm font-bold text-text-main">Immutable activity trail</h2></div><div className="mt-3 grid gap-2 sm:grid-cols-2 xl:grid-cols-3">{audit.slice(0, 9).map((event) => <div key={event.id} className="rounded-lg border border-borderSubtle bg-background/60 p-3"><p className="font-mono text-[10px] font-bold text-aurora">{event.action}</p><p className="mt-1 truncate text-xs text-text-main">{event.actor_id}</p><p className="mt-1 font-mono text-[10px] text-text-subtle">{new Date(event.occurred_at).toLocaleString("en-CA")}</p></div>)}</div></section>}
      </section></div>
    {message && <div role="status" className="fixed bottom-5 right-5 max-w-md rounded-xl border border-border bg-card px-4 py-3 text-xs text-text-main shadow-2xl">{message}</div>}
  </main>;
}
