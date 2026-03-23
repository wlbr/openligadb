package openligadb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://api.openligadb.de"

// Client is an HTTP client for the OpenLigaDB API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// NewClient creates a new OpenLigaDB client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func doGet[T any](ctx context.Context, c *Client, path string) (T, error) {
	var zero T

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return zero, fmt.Errorf("openligadb: creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("openligadb: executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("openligadb: unexpected status %d for %s", resp.StatusCode, path)
	}

	var result T
	if err := json.NewDecoder(io.LimitReader(resp.Body, 10<<20)).Decode(&result); err != nil {
		return zero, fmt.Errorf("openligadb: decoding response: %w", err)
	}
	return result, nil
}

// GetAvailableLeagues returns all available leagues.
func (c *Client) GetAvailableLeagues(ctx context.Context) ([]League, error) {
	return doGet[[]League](ctx, c, "/getavailableleagues")
}

// GetAvailableSports returns all available sports.
func (c *Client) GetAvailableSports(ctx context.Context) ([]Sport, error) {
	return doGet[[]Sport](ctx, c, "/getavailablesports")
}

// GetMatch returns a single match by its ID.
func (c *Client) GetMatch(ctx context.Context, matchID int) (Match, error) {
	return doGet[Match](ctx, c, fmt.Sprintf("/getmatchdata/%d", matchID))
}

// GetMatchesByLeagueSeason returns all matches for a league and season.
func (c *Client) GetMatchesByLeagueSeason(ctx context.Context, leagueShortcut string, leagueSeason int) ([]Match, error) {
	return doGet[[]Match](ctx, c, fmt.Sprintf("/getmatchdata/%s/%d", url.PathEscape(leagueShortcut), leagueSeason))
}

// GetMatchesByLeagueSeasonGroup returns matches for a league, season and matchday.
func (c *Client) GetMatchesByLeagueSeasonGroup(ctx context.Context, leagueShortcut string, leagueSeason, groupOrderID int) ([]Match, error) {
	return doGet[[]Match](ctx, c, fmt.Sprintf("/getmatchdata/%s/%d/%d", url.PathEscape(leagueShortcut), leagueSeason, groupOrderID))
}

// GetMatchesByLeagueSeasonTeam returns matches for a league, season filtered by team name.
func (c *Client) GetMatchesByLeagueSeasonTeam(ctx context.Context, leagueShortcut string, leagueSeason int, teamFilter string) ([]Match, error) {
	return doGet[[]Match](ctx, c, fmt.Sprintf("/getmatchdata/%s/%d/%s", url.PathEscape(leagueShortcut), leagueSeason, url.PathEscape(teamFilter)))
}

// GetMatchesByTeamIDs returns all matches between two teams.
func (c *Client) GetMatchesByTeamIDs(ctx context.Context, teamID1, teamID2 int) ([]Match, error) {
	return doGet[[]Match](ctx, c, fmt.Sprintf("/getmatchdata/%d/%d", teamID1, teamID2))
}

// GetLastChangeDate returns the last change date for a league/season/matchday.
func (c *Client) GetLastChangeDate(ctx context.Context, leagueShortcut string, leagueSeason, groupOrderID int) (time.Time, error) {
	var t time.Time
	s, err := doGet[string](ctx, c, fmt.Sprintf("/getlastchangedate/%s/%d/%d", url.PathEscape(leagueShortcut), leagueSeason, groupOrderID))
	if err != nil {
		return t, err
	}
	return time.Parse(time.RFC3339, s)
}

// GetNextMatchByLeagueTeam returns the next match for a league and team.
func (c *Client) GetNextMatchByLeagueTeam(ctx context.Context, leagueID, teamID int) (Match, error) {
	return doGet[Match](ctx, c, fmt.Sprintf("/getnextmatchbyleagueteam/%d/%d", leagueID, teamID))
}

// GetNextMatchByLeagueShortcut returns the next match for a league shortcut.
func (c *Client) GetNextMatchByLeagueShortcut(ctx context.Context, leagueShortcut string) (Match, error) {
	return doGet[Match](ctx, c, fmt.Sprintf("/getnextmatchbyleagueshortcut/%s", url.PathEscape(leagueShortcut)))
}

// GetLastMatchByLeagueShortcut returns the last match for a league shortcut.
func (c *Client) GetLastMatchByLeagueShortcut(ctx context.Context, leagueShortcut string) (Match, error) {
	return doGet[Match](ctx, c, fmt.Sprintf("/getlastmatchbyleagueshortcut/%s", url.PathEscape(leagueShortcut)))
}

// GetLastMatchByLeagueTeam returns the last match for a league and team.
func (c *Client) GetLastMatchByLeagueTeam(ctx context.Context, leagueID, teamID int) (Match, error) {
	return doGet[Match](ctx, c, fmt.Sprintf("/getlastmatchbyleagueteam/%d/%d", leagueID, teamID))
}

// GetCurrentGroup returns the current group (matchday) for a league.
func (c *Client) GetCurrentGroup(ctx context.Context, leagueShortcut string) (Group, error) {
	return doGet[Group](ctx, c, fmt.Sprintf("/getcurrentgroup/%s", url.PathEscape(leagueShortcut)))
}

// GetResultInfos returns the configured result types for a league.
func (c *Client) GetResultInfos(ctx context.Context, leagueID int) (ResultInfo, error) {
	return doGet[ResultInfo](ctx, c, fmt.Sprintf("/getresultinfos/%d", leagueID))
}

// GetAvailableGroups returns all groups (matchdays) for a league and season.
func (c *Client) GetAvailableGroups(ctx context.Context, leagueShortcut string, leagueSeason int) ([]Group, error) {
	return doGet[[]Group](ctx, c, fmt.Sprintf("/getavailablegroups/%s/%d", url.PathEscape(leagueShortcut), leagueSeason))
}

// GetGoalGetters returns the top scorers for a league and season.
func (c *Client) GetGoalGetters(ctx context.Context, leagueShortcut string, leagueSeason int) ([]GoalGetter, error) {
	return doGet[[]GoalGetter](ctx, c, fmt.Sprintf("/getgoalgetters/%s/%d", url.PathEscape(leagueShortcut), leagueSeason))
}

// GetAvailableTeams returns all teams for a league and season.
func (c *Client) GetAvailableTeams(ctx context.Context, leagueShortcut string, leagueSeason int) ([]Team, error) {
	return doGet[[]Team](ctx, c, fmt.Sprintf("/getavailableteams/%s/%d", url.PathEscape(leagueShortcut), leagueSeason))
}

// GetBlTable returns the league table for a league and season.
func (c *Client) GetBlTable(ctx context.Context, leagueShortcut string, leagueSeason int) ([]BlTableTeam, error) {
	return doGet[[]BlTableTeam](ctx, c, fmt.Sprintf("/getbltable/%s/%d", url.PathEscape(leagueShortcut), leagueSeason))
}

// GetGroupTable returns the group table for a league and season.
func (c *Client) GetGroupTable(ctx context.Context, leagueShortcut string, leagueSeason int) ([]BlTableTeam, error) {
	return doGet[[]BlTableTeam](ctx, c, fmt.Sprintf("/getgrouptable/%s/%d", url.PathEscape(leagueShortcut), leagueSeason))
}

// GetMatchesByTeam returns matches for a team within a time range.
func (c *Client) GetMatchesByTeam(ctx context.Context, teamFilter string, weekCountPast, weekCountFuture int) ([]Match, error) {
	return doGet[[]Match](ctx, c, fmt.Sprintf("/getmatchesbyteam/%s/%d/%d", url.PathEscape(teamFilter), weekCountPast, weekCountFuture))
}

// GetMatchesByTeamID returns matches for a team ID within a time range.
func (c *Client) GetMatchesByTeamID(ctx context.Context, teamID, weekCountPast, weekCountFuture int) ([]Match, error) {
	return doGet[[]Match](ctx, c, fmt.Sprintf("/getmatchesbyteamid/%d/%d/%d", teamID, weekCountPast, weekCountFuture))
}
