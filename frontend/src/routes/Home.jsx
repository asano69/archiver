import { useParams } from "@solidjs/router";
import ArchiveViewer from "../components/archives/ArchiveViewer";

// Archive browser. The title list now lives in the sidebar (see
// components/layout/Sidebar.jsx); this route only renders the selected
// archive's saved HTML, filling the whole main area. Matches both "/"
// (params.id unset, shows a placeholder) and "/archives/:id" (see
// lib/router.jsx). TopBar and Sidebar render once in AppShell, not per
// route.
export default function Home() {
  const params = useParams();

  return (
    <div class="flex h-full min-h-0 w-full flex-1 flex-col">
      <ArchiveViewer id={params.id} />
    </div>
  );
}
