import {
  createSignal,
  onMount,
  onCleanup,
  For,
  Show,
} from "solid-js";
import { A } from "@solidjs/router";
import pb from "../../lib/pb";
import Loading from "../Loading";

const PAGE_SIZE = 30;

// Paginated list of archive titles, newest first, loaded page by page
// as the user scrolls (same IntersectionObserver pattern as
// components/notes/NotesList.jsx). Each row links to "/archives/:id",
// which ArchiveViewer (see routes/Home.jsx) reads to render that
// archive's saved HTML in the right pane.
//
// Props:
//   selectedId: id of the currently open archive, used to highlight
//     its row.
export default function ArchiveList(props) {
  const [archives, setArchives] = createSignal([]);
  const [page, setPage] = createSignal(0);
  const [hasMore, setHasMore] = createSignal(true);
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");

  // Set by the sentinel div's `ref` below; observed once mounted.
  let sentinel;
  let observer;

  const loadPage = async (pageNum) => {
    if (loading()) return;
    setLoading(true);
    setError("");
    try {
      const result = await pb
        .collection("archives")
        .getList(pageNum, PAGE_SIZE, { sort: "-created" });
      setArchives((prev) =>
        pageNum === 1 ? result.items : [...prev, ...result.items],
      );
      setHasMore(pageNum < result.totalPages);
      setPage(pageNum);
    } catch (err) {
      console.error(
        "[archives] failed to load:",
        `${err?.name}: ${err?.message}`,
        err?.stack,
      );
      setError(
        err?.data?.message || err?.message || "Failed to load archives.",
      );
    } finally {
      setLoading(false);
    }
  };

  onMount(() => {
    observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting && hasMore() && !loading()) {
        loadPage(page() + 1);
      }
    });
    if (sentinel) observer.observe(sentinel);
    loadPage(1);
  });

  onCleanup(() => observer?.disconnect());

  return (
    <div class="flex h-full min-h-0 flex-col">
      <Show when={error()}>
        <p class="p-2 text-sm text-[#dc3545]">{error()}</p>
      </Show>

      <div class="min-h-0 flex-1 overflow-y-auto">
        <For each={archives()}>
          {(archive) => (
            <A
              href={`/archives/${archive.id}`}
              class="block truncate border-b border-border px-3 py-2.5 text-sm text-text transition-colors hover:bg-hover-bg"
              classList={{ "bg-active-bg": archive.id === props.selectedId }}
            >
              {archive.title || "(untitled)"}
            </A>
          )}
        </For>

        <Show when={!loading() && archives().length === 0}>
          <p class="p-3 text-sm text-border">No archives yet.</p>
        </Show>

        <Show when={loading()}>
          <Loading />
        </Show>

        {/* Observed by IntersectionObserver to trigger the next page load. */}
        <div ref={sentinel} class="h-1" />
      </div>
    </div>
  );
}
