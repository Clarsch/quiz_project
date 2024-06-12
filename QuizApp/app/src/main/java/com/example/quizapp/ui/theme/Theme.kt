package com.example.quizapp.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val customColorScheme = lightColorScheme(
    background = QuizBlue,
    onBackground = Color.White,
    primary = QuizBlue,
    secondary = QuizYellow,

    )

@Composable
fun QuizAppTheme(
    content: @Composable () -> Unit
) {

    MaterialTheme(
        colorScheme = customColorScheme,
        typography = Typography,
        content = content
    )
}