// Las comunicaciones se basan en mensajes de JSON por el STDOUT
// Cada mensaje tiene un ID de operación, que especifica qué tipo de mensaje es. Los IDs son:
// - 0: output general, devuelto en operaciones como list. Contiene un campo "content" con la respuesta.
// - 1-9: errors, Contiene un campo "error" con el mensaje de error.
//   - 1: critical error, este error implica un cierre inmediato de la aplicación
//   - 2: regular error, este error implica errores menores que no necesariamente implican el cierre de la aplicación.
//   - 3-9: reserved for future use
//
// - 10-19: status
//   - 10: progress, Contiene los campos:
//       - "global"
//         - "total", total global steps
//         - "current", current global step
//	       - "name", name of the global step
//       - "partial", (can be null)
//         - "total", total partial steps
//         - "current", current partial step
//         - "details", details of the partial step
//   - 11: paused, returned when all works have been paused
//   - 12: resumed, returned when the resume request is being processed
//   - 13: cancel, returned when the cancel request is being processed
//   - 14-19: reserved for future use
//
// - 20-29: interactive, contiene los campos:
//     - "interaction_id", id of the interaction. It must be returned in the reply.
//     - "prompt", question asked
//     - "valid_responses", array of valid responses
//   - 20: error
//   - 21-29: reserved for future use

// Las comunicaciones se basan en mensajes de JSON por el STDIN
// Cada mensaje tiene un ID de operación, que especifica qué tipo de mensaje es. Los IDs son:
// - 10: print progress
// - 11: pause
// - 12: resumed
// - 13: cancel
// - 20: response for interactive output. It MUST contain:
//     - "interaction_id", id of the interaction returned by the output.
//     - "response", response. Must be one of the valid responses.

package io
