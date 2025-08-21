import pygame
import string
import os
from pathlib import Path

pygame.init()
screen = pygame.display.set_mode((1280, 720))
clock = pygame.time.Clock()
running = True

# monospaced font
font = pygame.font.SysFont('Consolas', 24)
char_width = font.size("A")[0]
char_height = font.size("A")[1]


char_map = [
]

render_input = False
input_position = (0,0)
camera_position = (0,0)
space_pressed = False
mouse_pressed = False

def screen_to_canvas(screen_coordinates: tuple):
    return (screen_coordinates[0] + camera_position[0],
            screen_coordinates[1] + camera_position[1])

def canvas_to_screen(canvas_coordinates: tuple):
    return (canvas_coordinates[0] - camera_position[0],
            canvas_coordinates[1] - camera_position[1])

def canvas_to_char(canvas_coordinates : tuple):
    x, y = canvas_coordinates

    return (x//char_width, y//char_height)

def char_to_canvas(char_coordinates : tuple):
    x, y = char_coordinates

    return (x * char_width, y * char_height)

def char_to_screen(char_coordinates : tuple):
    canvas = char_to_canvas(char_coordinates)

    return canvas_to_screen(canvas)

def render_char_at(char:str, char_coordinates : tuple):
    if (len(char) > 1):
        char = char[0]
    
    text_surface = font.render(char, False, (255,255,255))

    screen_coordinates = char_to_screen(char_coordinates)
    screen.blit(text_surface, screen_coordinates)

def render_input_at(char_coordinates : tuple):
    text_surface = font.render(" ", False, (255,255,255), (255,255,255))

    screen_coordinates = char_to_screen(char_coordinates)
    screen.blit(text_surface, screen_coordinates)

def add_char_at(char : str, char_coordinates: tuple):
    global char_map

    if (len(char) > 1):
        char = char[0]

    # ensure that there is no other char
    remove_char_at(char_coordinates)
    char_map.append({"char" : char, "position" : char_coordinates})

def get_char_at(char_coordinates : tuple):
    global char_map
    
    query = [ x for x in char_map if x["position"] == char_coordinates]

    return query[0]["char"] if len(query) > 0 else " "

def remove_char_at(char_coordinates: tuple):
    global char_map

    char_map = [x for x in char_map if x["position"] != char_coordinates]
    

def render_char_map():
    for obj in char_map:
        render_char_at(obj["char"], obj["position"])


def save(ascii):
    save_folder = Path("save")
    save_folder.mkdir(exist_ok=True)

    index=1
    while True:
        new_file = save_folder / f"art{index}.txt"

        if not new_file.exists():
            break
        index+=1

    new_file.write_text(ascii)

def create_ascii_from_map():

    def get_top_left_corner():
        min_x = float('inf')
        min_y = float('inf')

        for obj in char_map:
            x, y = obj["position"]
            if x < min_x:
                min_x = x
            if y < min_y:
                min_y = y
        
        # Se char_map estiver vazio, pode retornar 0,0
        if min_x == float('inf'): min_x = 0
        if min_y == float('inf'): min_y = 0

        return (min_x, min_y)

    def get_bottom_right_corner():
        max_x = float('-inf')
        max_y = float('-inf')

        for obj in char_map:
            x, y = obj["position"]
            if x > max_x:
                max_x = x
            if y > max_y:
                max_y = y

        # Se char_map estiver vazio, retorna 0,0
        if max_x == float('-inf'): max_x = 0
        if max_y == float('-inf'): max_y = 0

        return (max_x, max_y)
    
    top_left_corner = get_top_left_corner()
    bottom_right_corner = get_bottom_right_corner()

    ascii_str = "```\n"

    for y in range(top_left_corner[1],bottom_right_corner[1] + 1):
        for x in range(top_left_corner[0], bottom_right_corner[0] + 1):
            char_at_pos = get_char_at((x,y))
            print(f"char at: {char_at_pos}")
            ascii_str += char_at_pos
        ascii_str+="\n"

    ascii_str += "```"

    print(f"min : {top_left_corner[0]} {top_left_corner[0]} max: {bottom_right_corner[0]} {bottom_right_corner[1]}")
    print("SAVED: ")
    print(ascii_str)
    return ascii_str

while running:

    mouse_pos = pygame.mouse.get_pos()
    mouse_char_pos = canvas_to_char(screen_to_canvas(mouse_pos))

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
            break
            
        elif event.type == pygame.KEYDOWN:

            if (render_input):
                char = event.unicode
                print(f"Key pressed {char}")
                if event.key == pygame.K_SPACE:
                    space_pressed = True

                elif event.key == pygame.K_UP:
                    input_position = (input_position[0], input_position[1] -1)

                elif event.key == pygame.K_DOWN:
                    input_position = (input_position[0], input_position[1] +1)
                
                elif event.key == pygame.K_LEFT:
                    input_position = (input_position[0] -1, input_position[1])
                
                elif event.key == pygame.K_RIGHT:
                    input_position = (input_position[0] +1, input_position[1])

                elif event.key == pygame.K_BACKSPACE:
                    remove_char_at(input_position)
                    input_position = (input_position[0] - 1, input_position[1])

                elif (char and char in string.printable and char not in '\t\n\r\x0b\x0c'):
                    add_char_at(event.unicode,input_position)
                    input_position = (input_position[0] + 1, input_position[1])
                
                elif event.key == pygame.K_s and (event.mod & pygame.KMOD_CTRL):
                    save(
                        create_ascii_from_map()
                    )

        elif event.type == pygame.MOUSEMOTION:
            if space_pressed and mouse_pressed:
                x,y = event.rel
                camera_position = (camera_position[0] + -x, camera_position[1] + -y)
        
        elif event.type == pygame.KEYUP:
            if event.key == pygame.K_SPACE:
                space_pressed = False
        elif event.type == pygame.MOUSEBUTTONDOWN:
            render_input = True
            input_position = mouse_char_pos
            mouse_pressed = True
        elif event.type == pygame.MOUSEBUTTONUP:
            mouse_pressed= False

    screen_input_position = canvas_to_screen(char_to_canvas(input_position))

    # Verifica se cursor passou das extremidades da tela
    if screen_input_position[0] > screen.get_size()[0]:
        camera_position = (camera_position[0] + char_width, camera_position[1])

    if screen_input_position[1] > screen.get_size()[1]:
        camera_position = (camera_position[0], camera_position[1] + char_height)
    screen.fill("black")

    # render
    render_char_map()

    text_surface= font.render(f"pos : {camera_position[0]}, {camera_position[1]}", False, (255,255,255))
    screen.blit(text_surface, (0,0))

    if (render_input):
        render_input_at(input_position)

    # show
    pygame.display.flip()

    clock.tick(60)
pygame.quit()