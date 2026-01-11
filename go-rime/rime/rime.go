package rime

/*
// CGO Configuration
// =================
//
// Dependencies:
// 1. librime - RIME Input Method Engine library
//    - Installed via: brew install librime
//    - Library path: /opt/homebrew/opt/librime/lib/librime.dylib
//    - Header path: /opt/homebrew/opt/librime/include/
//    - Required headers: rime_api.h
//
// 2. Squirrel Input Method (optional, for data files)
//    - Data path: /Library/Input Methods/Squirrel.app/Contents/SharedSupport
//    - Provides: default.yaml, schema files, dictionaries
//    - Alternative: Download from https://github.com/rime/plum
//
// Build Configuration:
// - CFLAGS: Include path for RIME headers
// - LDFLAGS: Library path and linker flags for librime
//
// Note: This wrapper uses the new RIME API style (rime_get_api)
// which returns a function pointer structure, not the deprecated
// direct C function calls.

#cgo CFLAGS: -I/opt/homebrew/opt/librime/include
#cgo LDFLAGS: -L/opt/homebrew/opt/librime/lib -lrime
#include <rime_api.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <sys/stat.h>

// Global API pointer - initialized once via rime_get_api()
static RimeApi* api = NULL;

// Initialize the RIME API pointer
static void init_api() {
    if (!api) {
        api = rime_get_api();
    }
}

// Helper function to initialize RimeTraits structure
static inline void init_rime_traits(RimeTraits* traits) {
    traits->data_size = sizeof(RimeTraits) - sizeof(traits->data_size);
}

// Helper function to initialize RimeContext structure
static inline void init_rime_context(RimeContext* ctx) {
    ctx->data_size = sizeof(RimeContext) - sizeof(ctx->data_size);
}

// Helper function to initialize RimeStatus structure
static inline void init_rime_status(RimeStatus* status) {
    status->data_size = sizeof(RimeStatus) - sizeof(status->data_size);
}

// Helper function to initialize RimeCommit structure
static inline void init_rime_commit(RimeCommit* commit) {
    commit->data_size = sizeof(RimeCommit) - sizeof(commit->data_size);
}

// Wrapper functions using the new API style

static void rime_setup(RimeTraits* traits) {
    init_api();
    if (api && api->setup) {
        api->setup(traits);
    }
}

static void rime_initialize(RimeTraits* traits) {
    init_api();
    if (api && api->initialize) {
        api->initialize(traits);
    }
}

static void rime_finalize() {
    init_api();
    if (api && api->finalize) {
        api->finalize();
    }
}

static RimeSessionId rime_create_session() {
    init_api();
    if (api && api->create_session) {
        return api->create_session();
    }
    return 0;
}

static Bool rime_destroy_session(RimeSessionId session_id) {
    init_api();
    if (api && api->destroy_session) {
        return api->destroy_session(session_id);
    }
    return False;
}

static Bool rime_process_key(RimeSessionId session_id, int keycode, int mask) {
    init_api();
    if (api && api->process_key) {
        return api->process_key(session_id, keycode, mask);
    }
    return False;
}

static Bool rime_get_context(RimeSessionId session_id, RimeContext* context) {
    init_api();
    if (api && api->get_context) {
        return api->get_context(session_id, context);
    }
    return False;
}

static Bool rime_free_context(RimeContext* context) {
    init_api();
    if (api && api->free_context) {
        return api->free_context(context);
    }
    return False;
}

static Bool rime_get_status(RimeSessionId session_id, RimeStatus* status) {
    init_api();
    if (api && api->get_status) {
        return api->get_status(session_id, status);
    }
    return False;
}

static Bool rime_free_status(RimeStatus* status) {
    init_api();
    if (api && api->free_status) {
        return api->free_status(status);
    }
    return False;
}

static Bool rime_get_commit(RimeSessionId session_id, RimeCommit* commit) {
    init_api();
    if (api && api->get_commit) {
        return api->get_commit(session_id, commit);
    }
    return False;
}

static Bool rime_free_commit(RimeCommit* commit) {
    init_api();
    if (api && api->free_commit) {
        return api->free_commit(commit);
    }
    return False;
}

static Bool rime_commit_composition(RimeSessionId session_id) {
    init_api();
    if (api && api->commit_composition) {
        return api->commit_composition(session_id);
    }
    return False;
}

static void rime_clear_composition(RimeSessionId session_id) {
    init_api();
    if (api && api->clear_composition) {
        api->clear_composition(session_id);
    }
}

static void rime_set_option(RimeSessionId session_id, const char* option, Bool value) {
    init_api();
    if (api && api->set_option) {
        api->set_option(session_id, option, value);
    }
}

static Bool rime_get_option(RimeSessionId session_id, const char* option) {
    init_api();
    if (api && api->get_option) {
        return api->get_option(session_id, option);
    }
    return False;
}

static void rime_set_property(RimeSessionId session_id, const char* prop, const char* value) {
    init_api();
    if (api && api->set_property) {
        api->set_property(session_id, prop, value);
    }
}

static Bool rime_get_property(RimeSessionId session_id, const char* prop, char* value, size_t buffer_size) {
    init_api();
    if (api && api->get_property) {
        return api->get_property(session_id, prop, value, buffer_size);
    }
    return False;
}

static Bool rime_simulate_key_sequence(RimeSessionId session_id, const char* key_sequence) {
    init_api();
    if (api && api->simulate_key_sequence) {
        return api->simulate_key_sequence(session_id, key_sequence);
    }
    return False;
}

static Bool rime_select_candidate(RimeSessionId session_id, size_t index) {
    init_api();
    if (api && api->select_candidate) {
        return api->select_candidate(session_id, index);
    }
    return False;
}

static Bool rime_select_candidate_on_current_page(RimeSessionId session_id, size_t index) {
    init_api();
    if (api && api->select_candidate_on_current_page) {
        return api->select_candidate_on_current_page(session_id, index);
    }
    return False;
}

static Bool rime_is_maintenance_mode() {
    init_api();
    if (api && api->is_maintenance_mode) {
        return api->is_maintenance_mode();
    }
    return False;
}

static Bool rime_start_maintenance(Bool full_check) {
    init_api();
    if (api && api->start_maintenance) {
        return api->start_maintenance(full_check);
    }
    return False;
}

static void rime_join_maintenance_thread() {
    init_api();
    if (api && api->join_maintenance_thread) {
        api->join_maintenance_thread();
    }
}

// Helper to check if build directory exists
static int build_dir_exists(const char* user_data_dir) {
    if (!user_data_dir) return 0;
    
    char build_path[1024];
    snprintf(build_path, sizeof(build_path), "%s/build", user_data_dir);
    
    struct stat st;
    return (stat(build_path, &st) == 0 && S_ISDIR(st.st_mode));
}
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Rime holds the state of the RIME service
type Rime struct {
	initialized bool
}

// Traits represents RIME initialization traits
type Traits struct {
	SharedDataDir        string
	UserDataDir          string
	DistributionName     string
	DistributionCodeName string
	DistributionVersion  string
	AppName             string
}

// Session represents a RIME session
type Session struct {
	id C.RimeSessionId
}

// Candidate represents a candidate word
type Candidate struct {
	Text    string
	Comment string
}

// Menu represents the candidate menu
type Menu struct {
	PageSize                    int
	PageNo                     int
	IsLastPage                 bool
	HighlightedCandidateIndex int
	NumCandidates             int
	Candidates                []Candidate
	SelectKeys                string
}

// Context represents the input context
type Context struct {
	Composition         Composition
	Menu                Menu
	CommitTextPreview   string
	SelectLabels        []string
}

// Composition represents the composition area
type Composition struct {
	Length    int
	CursorPos int
	SelStart  int
	SelEnd    int
	Preedit   string
}

// Status represents the input method status
type Status struct {
	SchemaId      string
	SchemaName    string
	IsDisabled    bool
	IsComposing   bool
	IsAsciiMode   bool
	IsFullShape   bool
	IsSimplified  bool
	IsTraditional bool
	IsAsciiPunct  bool
}

// New creates a new RIME instance
func New(traits Traits) *Rime {
	rime := &Rime{}
	
	// Set up traits
	var cTraits C.RimeTraits
	C.init_rime_traits(&cTraits)
	
	if traits.SharedDataDir != "" {
		cTraits.shared_data_dir = C.CString(traits.SharedDataDir)
		defer C.free(unsafe.Pointer(cTraits.shared_data_dir))
	}
	if traits.UserDataDir != "" {
		cTraits.user_data_dir = C.CString(traits.UserDataDir)
		defer C.free(unsafe.Pointer(cTraits.user_data_dir))
	}
	if traits.DistributionName != "" {
		cTraits.distribution_name = C.CString(traits.DistributionName)
		defer C.free(unsafe.Pointer(cTraits.distribution_name))
	}
	if traits.DistributionCodeName != "" {
		cTraits.distribution_code_name = C.CString(traits.DistributionCodeName)
		defer C.free(unsafe.Pointer(cTraits.distribution_code_name))
	}
	if traits.DistributionVersion != "" {
		cTraits.distribution_version = C.CString(traits.DistributionVersion)
		defer C.free(unsafe.Pointer(cTraits.distribution_version))
	}
	if traits.AppName != "" {
		cTraits.app_name = C.CString(traits.AppName)
		defer C.free(unsafe.Pointer(cTraits.app_name))
	}
	
	// Initialize RIME
	C.rime_initialize(&cTraits)
	rime.initialized = true
	
	// Check if deployment is needed (build directory doesn't exist)
	// This happens on first run or when user data is freshly initialized
	needsDeploy := C.build_dir_exists(cTraits.user_data_dir) == 0
	if needsDeploy || C.rime_is_maintenance_mode() != 0 {
		// Start full maintenance/deploy to generate prebuilt files
		C.rime_start_maintenance(C.int(1)) // full_check = true
		C.rime_join_maintenance_thread()
	}
	
	// Set up finalizer to ensure cleanup
	runtime.SetFinalizer(rime, func(r *Rime) {
		if r.initialized {
			r.Shutdown()
		}
	})
	
	return rime
}

// Setup performs basic setup without initializing the service
func Setup(traits Traits) {
	var cTraits C.RimeTraits
	C.init_rime_traits(&cTraits)
	
	if traits.SharedDataDir != "" {
		cTraits.shared_data_dir = C.CString(traits.SharedDataDir)
		defer C.free(unsafe.Pointer(cTraits.shared_data_dir))
	}
	if traits.UserDataDir != "" {
		cTraits.user_data_dir = C.CString(traits.UserDataDir)
		defer C.free(unsafe.Pointer(cTraits.user_data_dir))
	}
	if traits.DistributionName != "" {
		cTraits.distribution_name = C.CString(traits.DistributionName)
		defer C.free(unsafe.Pointer(cTraits.distribution_name))
	}
	if traits.DistributionCodeName != "" {
		cTraits.distribution_code_name = C.CString(traits.DistributionCodeName)
		defer C.free(unsafe.Pointer(cTraits.distribution_code_name))
	}
	if traits.DistributionVersion != "" {
		cTraits.distribution_version = C.CString(traits.DistributionVersion)
		defer C.free(unsafe.Pointer(cTraits.distribution_version))
	}
	if traits.AppName != "" {
		cTraits.app_name = C.CString(traits.AppName)
		defer C.free(unsafe.Pointer(cTraits.app_name))
	}
	
	C.rime_setup(&cTraits)
}

// CreateSession creates a new RIME session
func (r *Rime) CreateSession() *Session {
	sessionID := C.rime_create_session()
	if sessionID == 0 {
		return nil
	}
	
	session := &Session{id: sessionID}
	
	// Set up finalizer to ensure session cleanup
	runtime.SetFinalizer(session, func(s *Session) {
		r.DestroySession(s)
	})
	
	return session
}

// DestroySession destroys a RIME session
func (r *Rime) DestroySession(session *Session) bool {
	if session == nil {
		return false
	}
	
	result := C.rime_destroy_session(session.id) != 0
	if result {
		session.id = 0
	}
	return result
}

// ProcessKey processes a keyboard event
func (r *Rime) ProcessKey(session *Session, keycode int, mask int) bool {
	if session == nil {
		return false
	}
	return C.rime_process_key(session.id, C.int(keycode), C.int(mask)) != 0
}

// GetContext retrieves the current input context
func (r *Rime) GetContext(session *Session) (*Context, error) {
	if session == nil {
		return nil, nil
	}
	
	var cCtx C.RimeContext
	C.init_rime_context(&cCtx)
	
	success := C.rime_get_context(session.id, &cCtx) != 0
	if !success {
		return nil, nil
	}
	
	defer C.rime_free_context(&cCtx)
	
	ctx := &Context{}
	
	// Fill composition
	if cCtx.composition.preedit != nil {
		ctx.Composition.Preedit = C.GoString(cCtx.composition.preedit)
	}
	ctx.Composition.Length = int(cCtx.composition.length)
	ctx.Composition.CursorPos = int(cCtx.composition.cursor_pos)
	ctx.Composition.SelStart = int(cCtx.composition.sel_start)
	ctx.Composition.SelEnd = int(cCtx.composition.sel_end)
	
	// Fill menu
	ctx.Menu.PageSize = int(cCtx.menu.page_size)
	ctx.Menu.PageNo = int(cCtx.menu.page_no)
	ctx.Menu.IsLastPage = cCtx.menu.is_last_page != 0
	ctx.Menu.HighlightedCandidateIndex = int(cCtx.menu.highlighted_candidate_index)
	ctx.Menu.NumCandidates = int(cCtx.menu.num_candidates)
	
	// Fill candidates
	for i := 0; i < int(cCtx.menu.num_candidates); i++ {
		candidate := C.RimeCandidate{}
		if int(cCtx.menu.num_candidates) > i {
			candidate = *(*C.RimeCandidate)(unsafe.Pointer(uintptr(unsafe.Pointer(cCtx.menu.candidates)) + 
				uintptr(i)*unsafe.Sizeof(candidate)))
		}
		
		cand := Candidate{}
		if candidate.text != nil {
			cand.Text = C.GoString(candidate.text)
		}
		if candidate.comment != nil {
			cand.Comment = C.GoString(candidate.comment)
		}
		
		ctx.Menu.Candidates = append(ctx.Menu.Candidates, cand)
	}
	
	if cCtx.menu.select_keys != nil {
		ctx.Menu.SelectKeys = C.GoString(cCtx.menu.select_keys)
	}
	
	// Commit text preview
	if cCtx.commit_text_preview != nil {
		ctx.CommitTextPreview = C.GoString(cCtx.commit_text_preview)
	}
	
	// Select labels
	if cCtx.select_labels != nil {
		for i := 0; i < int(cCtx.menu.page_size); i++ {
			label := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cCtx.select_labels)) + 
				uintptr(i)*unsafe.Sizeof((*C.char)(nil))))
			if label != nil {
				ctx.SelectLabels = append(ctx.SelectLabels, C.GoString(label))
			}
		}
	}
	
	return ctx, nil
}

// GetStatus retrieves the current input method status
func (r *Rime) GetStatus(session *Session) (*Status, error) {
	if session == nil {
		return nil, nil
	}
	
	var cStatus C.RimeStatus
	C.init_rime_status(&cStatus)
	
	success := C.rime_get_status(session.id, &cStatus) != 0
	if !success {
		return nil, nil
	}
	
	defer C.rime_free_status(&cStatus)
	
	status := &Status{}
	
	if cStatus.schema_id != nil {
		status.SchemaId = C.GoString(cStatus.schema_id)
	}
	if cStatus.schema_name != nil {
		status.SchemaName = C.GoString(cStatus.schema_name)
	}
	status.IsDisabled = cStatus.is_disabled != 0
	status.IsComposing = cStatus.is_composing != 0
	status.IsAsciiMode = cStatus.is_ascii_mode != 0
	status.IsFullShape = cStatus.is_full_shape != 0
	status.IsSimplified = cStatus.is_simplified != 0
	status.IsTraditional = cStatus.is_traditional != 0
	status.IsAsciiPunct = cStatus.is_ascii_punct != 0
	
	return status, nil
}

// GetCommit retrieves the commit text
func (r *Rime) GetCommit(session *Session) (string, error) {
	if session == nil {
		return "", nil
	}
	
	var cCommit C.RimeCommit
	C.init_rime_commit(&cCommit)
	
	success := C.rime_get_commit(session.id, &cCommit) != 0
	if !success {
		return "", nil
	}
	
	defer C.rime_free_commit(&cCommit)
	
	if cCommit.text != nil {
		return C.GoString(cCommit.text), nil
	}
	
	return "", nil
}

// CommitComposition commits the current composition
func (r *Rime) CommitComposition(session *Session) bool {
	if session == nil {
		return false
	}
	return C.rime_commit_composition(session.id) != 0
}

// ClearComposition clears the current composition
func (r *Rime) ClearComposition(session *Session) {
	if session != nil {
		C.rime_clear_composition(session.id)
	}
}

// SetOption sets a session option
func (r *Rime) SetOption(session *Session, option string, value bool) {
	if session == nil {
		return
	}
	cOption := C.CString(option)
	defer C.free(unsafe.Pointer(cOption))
	
	var cValue C.int
	if value {
		cValue = 1
	}
	C.rime_set_option(session.id, cOption, cValue)
}

// GetOption gets a session option
func (r *Rime) GetOption(session *Session, option string) bool {
	if session == nil {
		return false
	}
	cOption := C.CString(option)
	defer C.free(unsafe.Pointer(cOption))
	
	return C.rime_get_option(session.id, cOption) != 0
}

// SetProperty sets a session property
func (r *Rime) SetProperty(session *Session, prop string, value string) {
	if session == nil {
		return
	}
	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	
	C.rime_set_property(session.id, cProp, cValue)
}

// GetProperty gets a session property
func (r *Rime) GetProperty(session *Session, prop string) string {
	if session == nil {
		return ""
	}
	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))
	
	buffer := make([]byte, 256)
	cBuffer := (*C.char)(unsafe.Pointer(&buffer[0]))
	
	success := C.rime_get_property(session.id, cProp, cBuffer, C.size_t(len(buffer))) != 0
	if !success {
		return ""
	}
	
	return C.GoString(cBuffer)
}

// SimulateKeySequence simulates a sequence of key presses
func (r *Rime) SimulateKeySequence(session *Session, keySequence string) bool {
	if session == nil {
		return false
	}
	cKeySeq := C.CString(keySequence)
	defer C.free(unsafe.Pointer(cKeySeq))
	
	return C.rime_simulate_key_sequence(session.id, cKeySeq) != 0
}

// SelectCandidate selects a candidate at the given index
func (r *Rime) SelectCandidate(session *Session, index int) bool {
	if session == nil {
		return false
	}
	return C.rime_select_candidate(session.id, C.size_t(index)) != 0
}

// SelectCandidateOnCurrentPage selects a candidate on the current page
func (r *Rime) SelectCandidateOnCurrentPage(session *Session, index int) bool {
	if session == nil {
		return false
	}
	return C.rime_select_candidate_on_current_page(session.id, C.size_t(index)) != 0
}

// Shutdown shuts down the RIME service
func (r *Rime) Shutdown() {
	if r.initialized {
		C.rime_finalize()
		r.initialized = false
	}
}

// IsMaintenanceMode checks if RIME is in maintenance mode
func (r *Rime) IsMaintenanceMode() bool {
	return C.rime_is_maintenance_mode() != 0
}

// StartMaintenance starts maintenance mode
func (r *Rime) StartMaintenance(fullCheck bool) bool {
	var check C.int
	if fullCheck {
		check = 1
	}
	return C.rime_start_maintenance(check) != 0
}

// JoinMaintenanceThread waits for maintenance thread to finish
func (r *Rime) JoinMaintenanceThread() {
	C.rime_join_maintenance_thread()
}