BEGIN;

CREATE TABLE IF NOT EXISTS public.apps (
  id                      uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  organization_id         uuid NOT NULL,
  display_name            character varying(255) NOT NULL,
  slug                    character varying(255) NOT NULL,
  description             text NOT NULL DEFAULT '',
  canvas_id               uuid,
  code_storage_repo_id    character varying(255) NOT NULL DEFAULT '',
  code_storage_remote_url character varying(1024) NOT NULL DEFAULT '',
  default_branch          character varying(255) NOT NULL DEFAULT '',
  live_commit_sha         character varying(255) NOT NULL DEFAULT '',
  edit_session_branch     character varying(255),
  sync_status             character varying(32) NOT NULL DEFAULT 'ok',
  sync_error              text,
  created_by              uuid,
  created_at              timestamp with time zone NOT NULL DEFAULT now(),
  updated_at              timestamp with time zone NOT NULL DEFAULT now(),
  deleted_at              timestamp with time zone
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_apps_slug ON public.apps (slug) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_apps_organization_id ON public.apps (organization_id);
CREATE INDEX IF NOT EXISTS idx_apps_deleted_at ON public.apps (deleted_at);

COMMIT;
