import { createResource, Show } from "solid-js";
import pb from "../../lib/pb";
import Loading from "../Loading";

// Injects `<base target="_blank">` so every link in the archived page
// opens in a new tab. This has to be plain HTML with no JS: the sandbox
// below intentionally omits allow-scripts, so there's no way to attach
// a click handler inside the iframe, and no allow-same-origin means the
// parent can't reach into the iframe's document either.
function withBlankTargetLinks(html) {
  const baseTag = '<base target="_blank">';
  if (/<head[^>]*>/i.test(html)) {
    return html.replace(/<head[^>]*>/i, (tag) => `${tag}${baseTag}`);
  }
  if (/<html[^>]*>/i.test(html)) {
    return html.replace(/<html[^>]*>/i, (tag) => `${tag}<head>${baseTag}</head>`);
  }
  return `${baseTag}${html}`;
}

// Renders a single archive's saved HTML in an iframe, isolating its
// styles/scripts from the app shell instead of injecting the markup
// directly into the page.
//
// The HTML is fetched as text and set via `srcdoc` rather than pointing
// `src` straight at the PocketBase file URL: PocketBase serves uploaded
// files with "Content-Disposition: attachment" (to stop stored HTML
// from ever rendering on the PocketBase origin itself), which makes a
// sandboxed iframe's direct navigation to that URL get blocked as a
// download instead of rendering. fetch() isn't affected by that header,
// so it still returns the raw bytes.
//
// sandbox only grants allow-popups + allow-popups-to-escape-sandbox, so
// clicking a link (see withBlankTargetLinks above) can open a normal,
// unsandboxed tab. No allow-scripts, since archived pages are arbitrary
// third-party HTML (captured by the SingleFile extension, see
// internal/serve/handler.go) and letting them execute script isn't
// worth the risk for read-only viewing; and no allow-same-origin, so a
// srcdoc document gets its own unique/opaque origin instead of
// inheriting the app's -- keeping it unable to reach this app's cookies
// or auth token even if a script slipped through.
//
// Props:
//   id: archive record id to display, or undefined to show a
//     placeholder (nothing selected yet).
export default function ArchiveViewer(props) {
  const [archive] = createResource(
    () => props.id,
    (id) => pb.collection("archives").getOne(id),
  );

  const [html] = createResource(archive, async (record) => {
    const res = await fetch(pb.files.getURL(record, record.file));
    if (!res.ok) {
      throw new Error(`Failed to fetch archive file: ${res.status}`);
    }
    return res.text();
  });

  return (
    <div class="flex h-full min-h-0 w-full flex-col">
      <Show
        when={props.id}
        fallback={
          <div class="flex h-full items-center justify-center text-sm text-border">
            Select an archive to view it.
          </div>
        }
      >
        <Show when={!archive.loading && !html.loading} fallback={<Loading />}>
          <Show
            when={archive() && !html.error}
            fallback={
              <p class="p-3 text-sm text-[#dc3545]">
                {html.error ? "Failed to load archive content." : "Archive not found."}
              </p>
            }
          >
            <iframe
              srcdoc={withBlankTargetLinks(html())}
              title={archive().title || "Archive"}
              sandbox="allow-popups allow-popups-to-escape-sandbox"
              class="h-full w-full flex-1 rounded-md border border-border bg-bg"
            />
          </Show>
        </Show>
      </Show>
    </div>
  );
}
