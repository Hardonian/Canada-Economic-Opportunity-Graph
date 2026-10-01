import { useEffect, useState, useRef, useCallback } from 'react';

export interface LiveProjectEvent {
  id: string;
  project_id: string;
  event_type: string;
  event_date: string;
  title: string;
  description: string;
  evidence_id?: string;
}

export interface UseLiveEventsOptions {
  projectId?: string;
  eventType?: string;
  maxBuffer?: number;
  enabled?: boolean;
}

export function useLiveEvents(options: UseLiveEventsOptions = {}) {
  const { projectId, eventType, maxBuffer = 50, enabled = true } = options;
  const [events, setEvents] = useState<LiveProjectEvent[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [lastHeartbeat, setLastHeartbeat] = useState<Date | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  const clearEvents = useCallback(() => {
    setEvents([]);
  }, []);

  useEffect(() => {
    if (!enabled || typeof window === 'undefined') {
      return;
    }

    const params = new URLSearchParams();
    if (projectId) params.set('project_id', projectId);
    if (eventType) params.set('event_type', eventType);

    const qs = params.toString();
    const url = `/api/v1/stream/events${qs ? `?${qs}` : ''}`;

    let es: EventSource | null = null;
    let reconnectTimeout: NodeJS.Timeout | null = null;

    function connect() {
      try {
        es = new EventSource(url);
        eventSourceRef.current = es;

        es.addEventListener('connected', () => {
          setIsConnected(true);
          setLastHeartbeat(new Date());
        });

        es.addEventListener('project_event', (e: MessageEvent) => {
          try {
            const ev: LiveProjectEvent = JSON.parse(e.data);
            setEvents((prev) => [ev, ...prev].slice(0, maxBuffer));
            setLastHeartbeat(new Date());
          } catch {
            // ignore malformed payloads
          }
        });

        es.onerror = () => {
          setIsConnected(false);
          es?.close();
          // Attempt exponential backoff reconnection
          reconnectTimeout = setTimeout(connect, 3000);
        };
      } catch {
        setIsConnected(false);
      }
    }

    connect();

    return () => {
      if (reconnectTimeout) clearTimeout(reconnectTimeout);
      if (es) es.close();
      eventSourceRef.current = null;
      setIsConnected(false);
    };
  }, [projectId, eventType, maxBuffer, enabled]);

  return { events, isConnected, lastHeartbeat, clearEvents };
}
