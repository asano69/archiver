import { useParams } from "@solidjs/router";
import ArchiveList from "../components/archives/ArchiveList";
import ArchiveViewer from "../components/archives/ArchiveViewer";

// Archive browser: a fixed-height two-pane layout, titles on the left
// and the selected archive's saved HTML on the right. Matches both "/"
// (params.id unset, right pane shows a placeholder) and
// "/archives/:id" (see lib/router.jsx). TopBar and Sidebar render once
// in AppShell, not per route.
export default function Home() {
  const params = useParams();

  return (
    <div class="flex h-full min-h-0 w-full gap-4">
      <aside class="w-72 shrink-0 min-h-0 overflow-hidden rounded-md border border-border bg-bg">
        <ArchiveList selectedId={params.id} />
      </aside>
      <section class="min-h-0 flex-1">
        <ArchiveViewer id={params.id} />
      </section>
    </div>
  );
}
