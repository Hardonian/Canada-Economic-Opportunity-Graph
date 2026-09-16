use reqwest::blocking::Client;
use reqwest::header::{HeaderMap, HeaderValue, AUTHORIZATION, USER_AGENT};
use serde::de::DeserializeOwned;
use serde_json::Value;
use std::time::Duration;

use crate::error::{CogError, CogResult};
use crate::models::{
    CapitalStack, CegsEnvelope, Organization, Project, ProjectList, ScoreBundle,
    SignalList,
};

/// Default timeout for all blocking HTTP requests.
const DEFAULT_TIMEOUT_SECS: u64 = 30;

/// SDK version string embedded in the User-Agent header.
const SDK_VERSION: &str = concat!("cog-sdk-rust/", env!("CARGO_PKG_VERSION"));

/// Filter parameters for project list queries.
#[derive(Debug, Default, Clone)]
pub struct ProjectFilter {
    pub sector: Option<String>,
    pub province: Option<String>,
    pub stage: Option<String>,
    pub min_capex_cad: Option<i64>,
    pub search: Option<String>,
    pub limit: Option<i32>,
    pub offset: Option<i32>,
}

/// The main COG API client.
///
/// ```no_run
/// use cog_sdk::CogClient;
///
/// let client = CogClient::new("https://api.canadaopportunitygraph.ca").unwrap();
/// let projects = client.list_projects(Default::default()).unwrap();
/// println!("{} projects returned", projects.projects.len());
/// ```
#[derive(Debug, Clone)]
pub struct CogClient {
    base_url: String,
    http: Client,
}

impl CogClient {
    /// Construct a client pointing at `base_url`.
    ///
    /// Optionally supply an API key via the `COG_API_KEY` environment variable
    /// or pass it through [`CogClient::with_api_key`].
    pub fn new(base_url: &str) -> CogResult<Self> {
        let base_url = base_url.trim_end_matches('/').to_string();
        if base_url.is_empty() {
            return Err(CogError::InvalidBaseUrl("base URL must not be empty".into()));
        }

        let mut default_headers = HeaderMap::new();
        default_headers.insert(
            USER_AGENT,
            HeaderValue::from_static(SDK_VERSION),
        );

        // If COG_API_KEY is set, pre-attach it as a Bearer token.
        if let Ok(key) = std::env::var("COG_API_KEY") {
            if let Ok(val) = HeaderValue::from_str(&format!("Bearer {key}")) {
                default_headers.insert(AUTHORIZATION, val);
            }
        }

        let http = Client::builder()
            .default_headers(default_headers)
            .timeout(Duration::from_secs(DEFAULT_TIMEOUT_SECS))
            .build()
            .map_err(CogError::Http)?;

        Ok(Self { base_url, http })
    }

    /// Return a new client with the given API key pre-attached.
    pub fn with_api_key(mut self, api_key: &str) -> CogResult<Self> {
        let val = HeaderValue::from_str(&format!("Bearer {api_key}"))
            .map_err(|_| CogError::InvalidBaseUrl("invalid API key characters".into()))?;

        // Rebuild client with updated headers.
        let mut headers = HeaderMap::new();
        headers.insert(USER_AGENT, HeaderValue::from_static(SDK_VERSION));
        headers.insert(AUTHORIZATION, val);

        self.http = Client::builder()
            .default_headers(headers)
            .timeout(Duration::from_secs(DEFAULT_TIMEOUT_SECS))
            .build()
            .map_err(CogError::Http)?;
        Ok(self)
    }

    // ── Projects ────────────────────────────────────────────────────────────

    /// List projects, optionally filtered.
    pub fn list_projects(&self, filter: ProjectFilter) -> CogResult<ProjectList> {
        let mut params: Vec<(&str, String)> = Vec::new();
        if let Some(v) = &filter.sector {
            params.push(("sector", v.clone()));
        }
        if let Some(v) = &filter.province {
            params.push(("province", v.clone()));
        }
        if let Some(v) = &filter.stage {
            params.push(("stage", v.clone()));
        }
        if let Some(v) = filter.min_capex_cad {
            params.push(("min_capex_cad", v.to_string()));
        }
        if let Some(v) = &filter.search {
            params.push(("search", v.clone()));
        }
        if let Some(v) = filter.limit {
            params.push(("limit", v.to_string()));
        }
        if let Some(v) = filter.offset {
            params.push(("offset", v.to_string()));
        }
        self.get("/api/v1/projects", &params)
    }

    /// Retrieve a single project by ID or slug.
    pub fn get_project(&self, id: &str) -> CogResult<Project> {
        self.get(&format!("/api/v1/projects/{id}"), &[])
    }

    /// Retrieve the latest scores for a project.
    pub fn get_project_scores(&self, project_id: &str) -> CogResult<ScoreBundle> {
        self.get(&format!("/api/v1/projects/{project_id}/scores"), &[])
    }

    /// Retrieve the resolved capital stack for a project.
    pub fn get_capital_stack(&self, project_id: &str) -> CogResult<CapitalStack> {
        self.get(&format!("/api/v1/projects/{project_id}/capital-stack"), &[])
    }

    // ── Signals ─────────────────────────────────────────────────────────────

    /// List recent momentum signals.
    ///
    /// `since_hours` — how far back to look (default 24 if `None`).
    pub fn list_signals(&self, since_hours: Option<u32>, limit: Option<i32>) -> CogResult<SignalList> {
        let mut params: Vec<(&str, String)> = Vec::new();
        if let Some(h) = since_hours {
            params.push(("since", format!("{h}h")));
        }
        if let Some(l) = limit {
            params.push(("limit", l.to_string()));
        }
        self.get("/api/v1/signals", &params)
    }

    // ── Organizations ────────────────────────────────────────────────────────

    /// List organizations.
    pub fn list_organizations(&self, limit: Option<i32>) -> CogResult<Vec<Organization>> {
        let mut params: Vec<(&str, String)> = Vec::new();
        if let Some(l) = limit {
            params.push(("limit", l.to_string()));
        }
        self.get("/api/v1/entities", &params)
    }

    // ── CEGS ────────────────────────────────────────────────────────────────

    /// Retrieve a project as a canonical CEGS 1.0 envelope.
    pub fn get_cegs_project(&self, project_id: &str) -> CogResult<CegsEnvelope> {
        self.get(&format!("/api/v1/cegs/projects/{project_id}"), &[])
    }

    // ── GraphQL ─────────────────────────────────────────────────────────────

    /// Execute a raw GraphQL query against the `/api/v1/graphql` endpoint.
    ///
    /// Returns the raw `data` object as a `serde_json::Value`.
    pub fn graphql(&self, query: &str, variables: Option<Value>) -> CogResult<Value> {
        let url = format!("{}/api/v1/graphql", self.base_url);
        let mut body = serde_json::json!({ "query": query });
        if let Some(vars) = variables {
            body["variables"] = vars;
        }

        let resp = self
            .http
            .post(&url)
            .header("Content-Type", "application/json")
            .json(&body)
            .send()?;

        let status = resp.status().as_u16();
        let text = resp.text()?;
        let mut parsed: Value = serde_json::from_str(&text)?;

        if let Some(errors) = parsed.get("errors") {
            if !errors.as_array().map_or(true, |a| a.is_empty()) {
                let msg = errors[0]["message"].as_str().unwrap_or("unknown").to_string();
                return Err(CogError::Api {
                    status,
                    code: "graphql_error".into(),
                    message: msg,
                });
            }
        }
        Ok(parsed["data"].take())
    }

    // ── Strategic Sovereignty Pillars ────────────────────────────────────────

    /// Get Pillar A: Critical Minerals taxonomy, refining clusters, and domestic retention statistics.
    pub fn get_critical_minerals(&self) -> CogResult<Value> {
        self.get("/api/v1/critical-minerals", &[])
    }

    /// Get Pillar B: Indigenous sovereign co-investment overview.
    pub fn get_indigenous_overview(&self) -> CogResult<Value> {
        self.get("/api/v1/indigenous", &[])
    }

    /// Simulate Pillar B: $5B Federal Indigenous Loan Guarantee Program syndication.
    pub fn simulate_indigenous_loan_guarantee(
        &self,
        capex_cad: Option<f64>,
        equity_pct: Option<f64>,
    ) -> CogResult<Value> {
        let mut params = Vec::new();
        if let Some(c) = capex_cad {
            params.push(("capex_cad", c.to_string()));
        }
        if let Some(e) = equity_pct {
            params.push(("indigenous_equity_pct", e.to_string()));
        }
        self.get("/api/v1/indigenous/loan-guarantee-sim", &params)
    }

    /// Get Pillar C: Internal Canadian trade barrier friction model.
    pub fn get_trade_friction(&self) -> CogResult<Value> {
        self.get("/api/v1/trade/friction", &[])
    }

    /// Get Pillar D: Clean Baseload and Sovereign AI Compute (FLOPs/MW metric).
    pub fn get_compute_sovereignty(&self) -> CogResult<Value> {
        self.get("/api/v1/compute/sovereignty", &[])
    }

    /// Query projects within a spatial bounding box.
    pub fn list_projects_spatial(
        &self,
        min_lat: f64,
        max_lat: f64,
        min_lng: f64,
        max_lng: f64,
        limit: Option<i32>,
    ) -> CogResult<Value> {
        let mut params = vec![
            ("min_lat", min_lat.to_string()),
            ("max_lat", max_lat.to_string()),
            ("min_lng", min_lng.to_string()),
            ("max_lng", max_lng.to_string()),
        ];
        if let Some(l) = limit {
            params.push(("limit", l.to_string()));
        }
        self.get("/api/v1/spatial/projects", &params)
    }

    // ── internal ─────────────────────────────────────────────────────────────

    fn get<T: DeserializeOwned>(&self, path: &str, params: &[(&str, String)]) -> CogResult<T> {
        let url = format!("{}{}", self.base_url, path);
        let mut req = self.http.get(&url);
        if !params.is_empty() {
            req = req.query(params);
        }
        let resp = req.send()?;
        let status = resp.status();
        let text = resp.text()?;

        if status.is_success() {
            let value: T = serde_json::from_str(&text)?;
            return Ok(value);
        }

        // Try to parse structured API error.
        if let Ok(parsed) = serde_json::from_str::<Value>(&text) {
            let code = parsed["error"].as_str().unwrap_or("api_error").to_string();
            let message = parsed["message"].as_str().unwrap_or(&text).to_string();
            return Err(CogError::Api {
                status: status.as_u16(),
                code,
                message,
            });
        }

        Err(CogError::Api {
            status: status.as_u16(),
            code: "api_error".into(),
            message: text,
        })
    }
}
