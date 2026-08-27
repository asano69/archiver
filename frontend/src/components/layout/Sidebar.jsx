import { Show } from "solid-js";
import { useParams } from "@solidjs/router";
import ArchiveList from "../archives/ArchiveList";

// Visibility is fully controlled by `open`; no separate desktop/mobile
// behavior.
export default function Sidebar(props) {
  // Highlights the currently open archive. Read directly from the route
  // (see routes/Home.jsx) instead of being passed down as a prop, since
  // Sidebar is rendered inside the router's context (see lib/router.jsx)
  // and can call useParams() itself.
  const params = useParams();

  return (
    <Show when={props.open}>
      {/* Backdrop only exists on mobile: clicking it closes the overlay,
          and it visually separates the sidebar from the content behind it. */}
      <Show when={props.isMobile}>
        <div
          class="absolute inset-0 z-20 bg-black/40"
          onClick={props.onClose}
        />
      </Show>
      <aside
        classList={{
          "absolute inset-y-0 left-0 z-30 shadow-popover": props.isMobile,
        }}
        class="flex h-full min-h-0 w-64 flex-col border-r border-border bg-bg"
      >
        <ArchiveList selectedId={params.id} />
      </aside>
    </Show>
  );
}
