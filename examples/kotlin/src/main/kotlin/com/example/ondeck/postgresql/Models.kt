// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package com.example.ondeck.postgresql

import java.time.LocalDateTime

// Venues can be either open or closed
enum class Status(val value: String) {
  OPEN("op!en"),
  CLOSED("clo@sed");

  companion object {
    private val map = Status.values().associateBy(Status::value)
    fun lookup(value: String) = map[value]
  }
}

data class City (
  val slug: String,
  val name: String
)

// Venues are places where muisc happens
data class Venue (
  val id: Int,
  val status: Status,
  val statuses: List<Status>,
  // This value appears in public URLs
  val slug: String,
  val name: String,
  val city: String,
  val spotifyPlaylist: String,
  val songkickId: String?,
  val tags: List<String>,
  val createdAt: LocalDateTime
)

